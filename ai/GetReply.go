package ai

import (
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
	resp := SendReq(aiModel, Msgs)
	if len(resp.Choices) == 0 {
		return ""
	}
	text := resp.Choices[0].Msg.Content
	loger.Loger.Info("[Ai]Ai说：", zap.String("text", text), zap.Int("本次消耗token", resp.Usage.TotalToken))
	return text
}
