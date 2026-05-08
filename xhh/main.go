package xhh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"xhhrobot/loger"
	"xhhrobot/pg"

	"go.uber.org/zap"
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
	ctx := context.Background()
	fmt.Println("[XHH]检查@")
	var offset int
	nomore := "false"
	other := fmt.Sprintf("?message_type=16&offset=%v&limit=20&no_more=%s", offset, nomore)
	resp := SendReq("GET", "/bbs/app/user/message", nil, other)
	var data Respo
	defer ctx.Done()
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
		_, err := pg.Conn.Exec(ctx, "INSERT INTO at (msg_id,comment_a_id,comment_root_id,link_id,user_a_id,comment_text,reply) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (msg_id) DO NOTHING", v.MsgID, v.CommentID, v.RootCommentID, v.LinkID, v.UserID, v.CommentText, false)
		if err != nil {
			loger.Loger.Info("[XHH]PsqlError", zap.Error(err))
			return
		}
	}
}

func AutoReply() {
	ctx := context.Background()
	row := pg.Conn.QueryRow(ctx, "SELECT link_id,comment_a_id,comment_root_id,comment_text FROM at WHERE reply=false LIMIT 1")
	var linkID, commentID, rootID int
	var text string
	row.Scan(&linkID, &commentID, &rootID, &text)
	if commentID != 0 {
		loger.Loger.Info(fmt.Sprintf("[XHH]正在回复[%v]%s", commentID, text))
		Reply("Ask Grok is currently available to Premium and Premium+ subscribers only. Subscribe to unlock this feature: x.com/i/premium_sign…", strconv.Itoa(linkID), strconv.Itoa(commentID), strconv.Itoa(rootID), "")
		pg.Conn.Exec(ctx, "UPDATE at SET reply=$1 WHERE comment_a_id=$2", true, commentID)
	} else {
		fmt.Println("[XHH]无事可做")
	}
	defer ctx.Done()
}
