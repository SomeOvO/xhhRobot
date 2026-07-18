package xhh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"sync"
	"time"
	"xhhrobot/db"
	"xhhrobot/loger"

	"go.uber.org/zap"
)

var lock = &sync.Mutex{}

func Reply(text, link_id, reply_id, root_id, iscy string) (isok bool) {
	lock.Lock()
	defer lock.Unlock()
	Path := "/bbs/app/comment/create"
	Body := fmt.Sprintf("is_cy=%s&link_id=%s&reply_id=%s&root_id=%s&text=%s", iscy, link_id, reply_id, root_id, url.QueryEscape(text))
	resp := SendReqReply("POST", Path, bytes.NewReader([]byte(Body)), "")
	if resp == nil {
		loger.Loger.Error("[XHH]链接发送失败了")
		return
	}
	defer resp.Body.Close()
	var resps struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		loger.Loger.Error("[XHH]无法解析Body", zap.Error(err))
		return false
	}
	err = json.Unmarshal(data, &resps)
	if err != nil {
		loger.Loger.Error("[XHH]无法反序列化", zap.String("body", string(data)), zap.Error(err))
		return false
	}
	if resps.Status != "ok" {
		if resps.Status == "failed" {
			CommentID, err := strconv.Atoi(reply_id)
			if err != nil {
				loger.Loger.Fatal("[XHH]不可能发生的事", zap.Error(err), zap.Any("info", resps), zap.Any("errs", reply_id))
			}
			db.Replyed(CommentID)
			loger.Loger.Info("[XHH]因为无法评论，所以已标记为完成", zap.Any("Resp", resps))
			time.Sleep(5 * time.Second)
			return true
		}
		if resps.Msg == "评论已被删除" {
			time.Sleep(5 * time.Second)
			return true
		}
		loger.Loger.Error("[XHH]评论发送失败", zap.Any("info", resps))
		return false
	}
	time.Sleep(5 * time.Second)
	return true
}
