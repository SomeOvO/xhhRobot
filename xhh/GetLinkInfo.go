package xhh

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"xhhrobot/ai"
	"xhhrobot/db"
	"xhhrobot/loger"

	"go.uber.org/zap"
)

type LinkInfoS struct {
	Msg    string `json:"msg"`
	Result struct {
		Link struct {
			Title  string      `json:"title"`
			Text   string      `json:"text"`
			Topics []ai.Topics `json:"topics"`
			Tags   []ai.Tags   `json:"hashtags"`
		} `json:"link"`
	} `json:"result"`
	Stat string `json:"status"`
}
type TextDetail struct {
	Text string `json:"text"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

func GetLinkInfo(LinkID int, CommentID int) (Contents []ai.Content, Topics []ai.Topics, Tags []ai.Tags) {
	resp := SendReq("GET", "/bbs/app/link/tree", nil, "?h_src&link_id="+strconv.Itoa(LinkID))
	if resp == nil {
		loger.Loger.Error("[XHH]获取帖子详情请求失败")
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		loger.Loger.Error("[XHH]无法读取响应体", zap.Error(err))
		return
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		loger.Loger.Error("[XHH]获取帖子详情失败", zap.Int("status", resp.StatusCode), zap.String("raw", string(data)))
		return
	}
	var RespS LinkInfoS

	err = json.Unmarshal(data, &RespS)
	if err != nil {
		loger.Loger.Error("[XHH]反序列化失败", zap.Error(err), zap.Any("data", string(data)))
		return
	}
	if RespS.Stat != "ok" {
		if RespS.Stat == "failed" {
			db.Replyed(CommentID)
			loger.Loger.Warn("[XHH]原帖无法被查看，所以已标记为完成")
			return
		}
		loger.Loger.Error("[XHH]返回了错误的内容", zap.Any("info", RespS), zap.Any("data", string(data)))
		return
	}
	var Content []TextDetail

	err = json.Unmarshal([]byte(RespS.Result.Link.Text), &Content)
	if err != nil {
		loger.Loger.Error("[XHH]无法解析内容", zap.Error(err))
		return
	}
	var Title ai.Content
	Title.Text = "以下是帖子内容：\n标题：" + RespS.Result.Link.Title
	Title.Type = "text"
	Contents = append(Contents, Title)
	for _, v := range Content {
		var content ai.Content
		if v.Type == "html" {
			content.Type = "text"
			content.Text = v.Text
			Contents = append(Contents, content)
			continue
		}
		if v.Type != "text" {
			content.Type = "image_url"
			content.ImgUrl.Url = v.Url
			Contents = append(Contents, content)
			continue
		}
		content.Type = "text"
		content.Text = v.Text
		Contents = append(Contents, content)
	}
	return Contents, RespS.Result.Link.Topics, RespS.Result.Link.Tags
}
