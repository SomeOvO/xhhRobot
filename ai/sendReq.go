package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"xhhrobot/config"
	"xhhrobot/loger"

	"go.uber.org/zap"
)

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BodyStruct struct {
	Model string     `json:"model"`
	Msgs  []Messages `json:"messages"`
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

func SendReq(Model string, Msg []Messages) (Jresp respStruct) {
	if Model == "" {
		loger.Loger.Fatal("[Ai]请确保配置文件中的模型是存在的")
	}
	var body BodyStruct
	body.Model = Model
	body.Msgs = Msg
	cfg := config.ConfigStruct.Ai
	reqbody, err := json.Marshal(body)
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

	Dresp, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(Dresp, &Jresp)
	if err != nil {
		loger.Loger.Error("[Ai]无法反序列化JSON", zap.Error(err), zap.String("body", string(Dresp)))
		return
	}

	return Jresp

}
