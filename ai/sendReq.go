package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"xhhrobot/config"
	"xhhrobot/loger"

	"go.uber.org/zap"
)

type Content struct {
	Type   string `json:"type"`
	ImgUrl struct {
		Url string `json:"url"`
	} `json:"image_url"`
	Text string `json:"text"`
}
type Messages[T []Content | string] struct {
	Role    string `json:"role"`
	Content T      `json:"content"`
}
type SysMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type BodyStruct struct {
	Model            string         `json:"model"`
	Msgs             []any          `json:"messages,omitempty"`
	Stream           bool           `json:"stream"`
	WebSearchOptions map[string]any `json:"web_search_options,omitempty"`
}

type responsesBodyStruct struct {
	Model      string              `json:"model"`
	Input      []responsesInputMsg `json:"input"`
	Tools      []responsesWebTool  `json:"tools,omitempty"`
	ToolChoice string              `json:"tool_choice,omitempty"`
	Stream     bool                `json:"stream"`
}

type responsesInputMsg struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type responsesInputContent struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageUrl string `json:"image_url,omitempty"`
}

type responsesWebTool struct {
	Type              string `json:"type"`
	SearchContextSize string `json:"search_context_size,omitempty"`
}

type choice struct {
	Index int `json:"index"`
	Msg   struct {
		Role    string `json:"role"`
		Content string `json:"content"`
		Reason  string `json:"reasoning_content"`
	} `json:"message"`
}

type respStruct struct {
	Choices []choice `json:"choices"`
	Usage   struct {
		TotalToken int `json:"total_tokens"`
	} `json:"usage"`
}

type rawMsg struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type responsesRespStruct struct {
	OutputText string `json:"output_text"`
	Output     []struct {
		Type    string `json:"type"`
		Role    string `json:"role"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Usage struct {
		TotalToken int `json:"total_tokens"`
	} `json:"usage"`
}

func SendReq(Model string, Msg []any) (Jresp respStruct) {
	if Model == "" {
		loger.Loger.Fatal("[Ai]请确保配置文件中的模型是存在的")
	}
	cfg := config.ConfigStruct.Ai
	reqbody, err := buildReqBody(Model, Msg)
	if err != nil {
		loger.Loger.Error("[AI]无法序列化JSON", zap.Error(err))
		return
	}
	req, err := http.NewRequest("POST", cfg.BaseUrl, bytes.NewReader(reqbody))
	if err != nil {
		loger.Loger.Error("[AI]无法创建请求", zap.Error(err))
		return
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		loger.Loger.Error("[AI]请求失败！", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	Dresp, err := io.ReadAll(resp.Body)
	if err != nil {
		loger.Loger.Error("[AI]无法读取响应", zap.Error(err))
		return
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		loger.Loger.Error("[Ai]Ai返回错误", zap.Int("status", resp.StatusCode), zap.String("body", string(Dresp)))
		return
	}
	if useResponsesAPI(cfg.BaseUrl) {
		Jresp, err = parseResponsesResp(Dresp)
	} else {
		err = json.Unmarshal(Dresp, &Jresp)
	}
	if err != nil {
		loger.Loger.Error("[Ai]无法反序列化JSON", zap.Error(err), zap.String("body", string(Dresp)))
		return
	}
	if len(Jresp.Choices) == 0 {
		loger.Loger.Error("[Ai]Ai返回错误", zap.Any("Resp", string(Dresp)))
		return
	}
	return Jresp

}

func buildReqBody(Model string, Msg []any) ([]byte, error) {
	cfg := config.ConfigStruct.Ai
	if useResponsesAPI(cfg.BaseUrl) {
		input, err := toResponsesInput(Msg)
		if err != nil {
			return nil, err
		}
		body := responsesBodyStruct{
			Model: Model,
			Input: input,
		}
		if cfg.WebSearch {
			tool := responsesWebTool{Type: "web_search"}
			if cfg.SearchContextSize != "" {
				tool.SearchContextSize = cfg.SearchContextSize
			}
			body.Tools = []responsesWebTool{tool}
			if cfg.ForceWebSearch {
				body.ToolChoice = "required"
			}
		}
		return json.Marshal(body)
	}

	body := BodyStruct{
		Model: Model,
		Msgs:  Msg,
	}
	if cfg.WebSearch {
		body.WebSearchOptions = map[string]any{}
		if strings.HasPrefix(strings.ToLower(Model), "gpt-5.4") {
			loger.Loger.Warn("[AI]Chat Completions 联网需要使用 gpt-5-search-api；gpt-5.4 联网请把 baseUrl 改为 /v1/responses")
		}
	}
	return json.Marshal(body)
}

func useResponsesAPI(baseUrl string) bool {
	return strings.Contains(strings.ToLower(baseUrl), "/responses")
}

func toResponsesInput(Msg []any) ([]responsesInputMsg, error) {
	input := make([]responsesInputMsg, 0, len(Msg))
	for _, msg := range Msg {
		data, err := json.Marshal(msg)
		if err != nil {
			return nil, err
		}
		var raw rawMsg
		if err := json.Unmarshal(data, &raw); err != nil {
			return nil, err
		}
		role := raw.Role
		if role == "system" {
			role = "developer"
		}
		content, err := toResponsesContent(raw.Content)
		if err != nil {
			return nil, err
		}
		input = append(input, responsesInputMsg{
			Role:    role,
			Content: content,
		})
	}
	return input, nil
}

func toResponsesContent(raw json.RawMessage) (any, error) {
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text, nil
	}

	var contents []Content
	if err := json.Unmarshal(raw, &contents); err != nil {
		return nil, err
	}
	responseContents := make([]responsesInputContent, 0, len(contents))
	for _, content := range contents {
		switch content.Type {
		case "text":
			responseContents = append(responseContents, responsesInputContent{
				Type: "input_text",
				Text: content.Text,
			})
		case "image_url":
			responseContents = append(responseContents, responsesInputContent{
				Type:     "input_image",
				ImageUrl: content.ImgUrl.Url,
			})
		default:
			if content.Text != "" {
				responseContents = append(responseContents, responsesInputContent{
					Type: "input_text",
					Text: content.Text,
				})
			}
		}
	}
	return responseContents, nil
}

func parseResponsesResp(data []byte) (respStruct, error) {
	var responsesResp responsesRespStruct
	if err := json.Unmarshal(data, &responsesResp); err != nil {
		return respStruct{}, err
	}

	text := responsesResp.OutputText
	if text == "" {
		var builder strings.Builder
		for _, output := range responsesResp.Output {
			if output.Type != "message" {
				continue
			}
			for _, content := range output.Content {
				if content.Type != "output_text" || content.Text == "" {
					continue
				}
				if builder.Len() > 0 {
					builder.WriteString("\n")
				}
				builder.WriteString(content.Text)
			}
		}
		text = builder.String()
	}
	if text == "" {
		return respStruct{}, nil
	}

	var resp respStruct
	resp.Choices = []choice{{Index: 0}}
	resp.Choices[0].Msg.Role = "assistant"
	resp.Choices[0].Msg.Content = text
	resp.Usage.TotalToken = responsesResp.Usage.TotalToken
	return resp, nil
}
