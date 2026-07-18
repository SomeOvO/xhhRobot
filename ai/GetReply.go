package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"xhhrobot/config"
	"xhhrobot/loger"

	"go.uber.org/zap"
)

type Topics struct {
	Name string `json:"name"`
}
type Tags struct {
	Name string `json:"name"`
}

func GetAiReply(Contents []Content, UserSay string, Topics []Topics, Tags []Tags) string {
	loger.Loger.Info("[Ai]正在询问Ai", zap.Any("Content", Contents))
	var SMsg Messages[string]
	var UMsg Messages[[]Content]
	var Msgs []any
	SMsg.Role = "system"
	cfg := config.ConfigStruct.Ai
	prompt := cfg.Prompt
	var topStr strings.Builder
	for _, v := range Topics {
		topStr.WriteString(v.Name)
	}
	prompt = strings.ReplaceAll(prompt, "?!top!?", topStr.String())
	var tagStr strings.Builder
	for _, v := range Tags {
		tagStr.WriteString(v.Name)
	}
	prompt = strings.ReplaceAll(prompt, "?!tag!?", tagStr.String())
	fmt.Println(prompt)
	SMsg.Content = prompt
	//用户
	UMsg.Role = "user"
	var UserContent Content
	UserContent.Text = "以上是帖子内容。" + UserSay
	UserContent.Type = "text"
	Contents = append(Contents, UserContent)
	UMsg.Content = Contents
	Msgs = append(Msgs, SMsg)
	Msgs = append(Msgs, UMsg)
	aiModel := config.ConfigStruct.Ai.Model

	maxRounds := config.ConfigStruct.Ai.MCP.MaxRounds
	if maxRounds <= 0 {
		maxRounds = 10
	}

	for round := 0; round < maxRounds; round++ {
		resp := SendReq(aiModel, Msgs)
		if len(resp.Choices) == 0 {
			loger.Loger.Error("[Ai]Ai返回错误", zap.Any("Resp", resp))
			return ""
		}
		c := resp.Choices[0]

		if c.FinishReason == "tool_calls" && len(c.Msg.ToolCalls) > 0 {
			Msgs = append(Msgs, AssistantToolMsg{
				Role:      "assistant",
				Content:   c.Msg.Content,
				ToolCalls: c.Msg.ToolCalls,
			})

			for _, tc := range c.Msg.ToolCalls {
				var args map[string]any
				if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
					loger.Loger.Error("[Ai]解析工具参数失败", zap.Error(err))
					args = nil
				}
				text := extractToolResult(mcpMgr.callTool(context.Background(), tc.Function.Name, args))
				loger.Loger.Info(
					"[Ai]工具调用",
					zap.String("tool", tc.Function.Name),
					zap.String("result", text),
				)
				Msgs = append(Msgs, ToolResultMsg{
					Role:       "tool",
					ToolCallID: tc.ID,
					Content:    text,
				})
			}
			continue
		}

		// 正常结束，返回文本
		text := c.Msg.Content
		loger.Loger.Info("[Ai]Ai说：", zap.String("text", text), zap.Int("本次消耗token", resp.Usage.TotalToken))
		return text
	}

	loger.Loger.Warn("[Ai]工具调用轮次耗尽", zap.Int("maxRounds", maxRounds))
	return ""
}
