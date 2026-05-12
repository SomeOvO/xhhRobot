package ai

import (
	"xhhrobot/config"
	"xhhrobot/loger"

	"go.uber.org/zap"
)

func GetAiReply(Contents []Content, UserSay string) string {
	loger.Loger.Info("[Ai]正在询问Ai", zap.Any("Content", Contents))
	var SMsg Messages[string]
	var UMsg Messages[[]Content]
	var Msgs []any
	SMsg.Role = "system"
	cfg := config.ConfigStruct.Ai
	SMsg.Content = cfg.Prompt
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
