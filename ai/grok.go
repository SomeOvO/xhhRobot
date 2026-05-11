package ai

import (
	"xhhrobot/config"
	"xhhrobot/loger"

	"go.uber.org/zap"
)

func Grok(Content, UserSend string) string {
	loger.Loger.Info("[Ai]正在询问Ai", zap.String("text", Content))
	var SMsg Messages
	var TMsg Messages
	var UMsg Messages
	var Msgs [3]Messages

	//系统提示词
	SMsg.Role = "system"
	cfg := config.ConfigStruct.Ai
	SMsg.Content = cfg.Prompt

	//帖子
	TMsg.Role = "user"
	TMsg.Content = "帖子内容：" + Content

	//用户
	UMsg.Role = "user"
	UMsg.Content = UserSend
	Msgs[0] = SMsg
	Msgs[1] = TMsg
	Msgs[2] = UMsg
	resp := SendReq(config.ConfigStruct.Ai.Model, Msgs[:])
	if len(resp.Choices) == 0 {
		return ""
	}
	text := resp.Choices[0].Msg.Content
	loger.Loger.Info("[Ai]Ai说：", zap.String("text", text), zap.Int("本次消耗token", resp.Usage.TotalToken))
	return text
}
