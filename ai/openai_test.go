package ai_test

import (
	"heybox/ai"
	"testing"
)

func OpenAi_ResponsesTest(T *testing.T) {
	Msg := ai.GetResponsesMsg()
	ai.AddResponsesTextMsg(&Msg, "你好")
	err, resp := ai.OpenAi_Responses(Msg)
	if err != nil {
		T.Fatal(err)
	}
	T.Log(resp)
	//还没写完
}
