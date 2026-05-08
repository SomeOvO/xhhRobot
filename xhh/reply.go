package xhh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"xhhrobot/loger"

	"go.uber.org/zap"
)

func Reply(text, link_id, reply_id, root_id, iscy string) (isok bool) {
	Path := "/bbs/app/comment/create"
	Body := fmt.Sprintf("is_cy=%s&link_id=%s&reply_id=%s&root_id=%s&text=%s", iscy, link_id, reply_id, root_id, url.QueryEscape(text))
	resp := SendReq("POST", Path, bytes.NewReader([]byte(Body)), "")
	var resps struct {
		Status string `json:"status"`
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		loger.Loger.Error("[XHH]无法解析Body", zap.Error(err))
		return false
	}
	err = json.Unmarshal(data, &resps)
	if err != nil {
		loger.Loger.Error("[XHH]无法反序列化", zap.Error(err))
		return false
	}
	if resps.Status != "ok" {
		loger.Loger.Error("[XHH]评论发送失败", zap.Any("info", resps))
		return false
	}
	return true
}
