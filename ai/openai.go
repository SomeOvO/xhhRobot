package ai

import (
	"context"
	"heybox/cfg"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func openai_client() *openai.Client {
	ApiKey := cfg.Config[cfg.CONFIG_AI_TOKEN]
	BaseUrl := cfg.Config[cfg.CONFIG_AI_BaseUrl]
	client := openai.NewClient(
		option.WithAPIKey(ApiKey.Value),
		option.WithBaseURL(BaseUrl.Value),
	)
	return &client
}

// 获取一个Msg列表
func GetResponsesMsg() responses.ResponseInputMessageContentListParam {
	Msg := responses.ResponseInputMessageContentListParam{}
	return Msg
}

// 增加文字至MsgList
func AddResponsesTextMsg(Msg *responses.ResponseInputMessageContentListParam, Content string) {
	*Msg = append(*Msg, responses.ResponseInputContentParamOfInputText(Content))
}

// 发送请求
func OpenAi_Responses(UserMsg responses.ResponseInputMessageContentListParam) (err error, resp *responses.Response) {
	client := openai_client()
	sysMsg := responses.ResponseInputItemParamOfMessage(
		responses.ResponseInputMessageContentListParam{},
		responses.EasyInputMessageRoleSystem,
	)
	userMsg := responses.ResponseInputItemParamOfMessage(
		UserMsg,
		responses.EasyInputMessageRoleUser,
	)
	Input := responses.ResponseNewParamsInputUnion{
		OfInputItemList: responses.ResponseInputParam{
			sysMsg, userMsg,
		},
	}
	resp, err = client.Responses.New(context.TODO(), responses.ResponseNewParams{
		Model: "",
		Input: Input,
	})
	if err != nil {
		return
	}
	return
}
