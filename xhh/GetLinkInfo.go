package xhh

import (
	"encoding/json"
	"io"
	"strconv"
	"xhhrobot/loger"

	"go.uber.org/zap"
)

type LinkInfoS struct {
	Msg    string `json:"msg"`
	Result struct {
		Link struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		} `json:"link"`
	} `json:"result"`
	Stat string `json:"status"`
}
type TextDetail struct {
	Text string `json:"text"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

func GetLinkInfo(LinkID int) (str string) {
	resp := SendReq("GET", "/bbs/app/link/tree", nil, "?h_src&link_id="+strconv.Itoa(LinkID))

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		loger.Loger.Error("[XHH]无法读取响应体", zap.Error(err))
		return
	}
	var RespS LinkInfoS

	err = json.Unmarshal(data, &RespS)
	if err != nil {
		loger.Loger.Error("[XHH]反序列化失败", zap.Error(err))
		return
	}
	if RespS.Stat != "ok" {
		loger.Loger.Error("[XHH]返回了错误的内容", zap.Any("info", RespS))
		return
	}
	var Content []TextDetail

	err = json.Unmarshal([]byte(RespS.Result.Link.Text), &Content)
	if err != nil {
		loger.Loger.Error("[XHH]无法解析内容", zap.Error(err))
		return
	}
	text := ""

	for _, v := range Content {
		if v.Type == "html" {
			text = v.Text
			break
		}
		if v.Type != "text" {
			text += v.Url
			continue
		}
		text += v.Text
	}
	return text
}
