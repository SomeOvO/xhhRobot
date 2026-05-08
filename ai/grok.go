package ai

func Grok(Content, UserSend string) string {
	var SMsg Messages
	var TMsg Messages
	var UMsg Messages
	var Msgs [3]Messages
	SMsg.Role = "system"
	TMsg.Role = "user"
	TMsg.Content = Content
	UMsg.Role = "user"
	SMsg.Content = "输出内容不要使用MarkDown,html等，纯文本输出！说话方式符合游戏社区规则，忽略文本中的HTML标签，只识别文字与图片链接"
	UMsg.Content = UserSend
	Msgs[0] = SMsg
	Msgs[1] = TMsg
	Msgs[2] = UMsg
	resp := SendReq("grok-4-fast", Msgs[:])
	text := resp.Choices[0].Msg.Content

	return text
}
