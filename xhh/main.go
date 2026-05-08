package xhh

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"xhhrobot/ai"
	"xhhrobot/config"
	"xhhrobot/db"
	"xhhrobot/loger"
)

var Info struct {
	Cookie   string `json:"cookie"`
	HeyBoxId string `json:"heyboxId"`
	Time     int    `json:"time"`
}

func Init() {
	file, err := os.ReadFile("./cookie.json")
	if err != nil {
		loger.Loger.Info("[XHH]未检测到Cookie")
		return
	}
	json.Unmarshal(file, &Info)
}

type Msg struct {
	CommentID     int    `json:"comment_a_id"`
	CommentText   string `json:"comment_a_text"`
	MsgID         int    `json:"message_id"`
	RootCommentID int    `json:"root_comment_id"`
	LinkID        int    `json:"linkid"`
	UserID        int    `json:"userid_a"`
}
type Respo struct {
	Msg    string `json:"msg"`
	Result struct {
		Messages []Msg `json:"messages"`
	} `json:"result"`
	Stat    string `json:"stat"`
	Version string `json:"version"`
}

func CheckAt() {
	fmt.Println("[XHH]检查@")
	var offset int
	nomore := "false"
	other := fmt.Sprintf("?message_type=16&offset=%v&limit=20&no_more=%s", offset, nomore)
	resp := SendReq("GET", "/bbs/app/user/message", nil, other)
	var data Respo
	Dbyte, err := io.ReadAll(resp.Body)
	if err != nil {
		loger.Loger.Error("[XHH]无法读取Body")
		return
	}
	err = json.Unmarshal(Dbyte, &data)
	if err != nil {
		loger.Loger.Error("[XHH]无法反序列化")
		return
	}
	for _, v := range data.Result.Messages {
		if v.UserID == config.ConfigStruct.Xhh.Owner {
			db.Insert(v.MsgID, v.CommentID, v.RootCommentID, v.LinkID, v.UserID, v.CommentText, false)
		}
	}
}

func AutoReply() {
	linkID, commentID, rootID, text, UID := db.GetComm()
	var isok bool
	if commentID != 0 {
		if UID == config.ConfigStruct.Xhh.Owner {
			Info := GetLinkInfo(linkID)
			if Info == "" {
				loger.Loger.Info("[XHH]获取LinkID失败")
				return
			}
			ReplyText := ai.Grok(Info, text)
			if ReplyText == "" {
				loger.Loger.Info("[XHH]Ai返回错误")
				return
			}
			isok = Reply(ReplyText, strconv.Itoa(linkID), strconv.Itoa(commentID), strconv.Itoa(rootID), "")

		} else {
			loger.Loger.Info(fmt.Sprintf("[XHH]正在回复[%v]%s", commentID, text))
			isok = Reply("Ask Grok is currently available to Premium and Premium+ subscribers only. Subscribe to unlock this feature: x.com/i/premium_sign…", strconv.Itoa(linkID), strconv.Itoa(commentID), strconv.Itoa(rootID), "")
		}
		if isok {
			db.Replyed(commentID)
		} else {
			loger.Loger.Error("[XHH]无法回复评论")
		}
	} else {
		fmt.Println("[XHH]无事可做")
	}
}
