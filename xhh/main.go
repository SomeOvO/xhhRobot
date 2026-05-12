package xhh

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
	"xhhrobot/ai"
	"xhhrobot/config"
	"xhhrobot/db"
	"xhhrobot/loger"

	"go.uber.org/zap"
)

var Info struct {
	Cookie   string `json:"cookie"`
	HeyBoxId string `json:"heyboxId"`
	Time     int    `json:"time"`
}
var CheckTime int
var ReplyTime int

func Init() {
	file, err := os.ReadFile("./cookie.json")
	if err != nil {
		loger.Loger.Info("[XHH]未检测到Cookie")
		return
	}
	CheckTime = config.ConfigStruct.Xhh.CheckTime
	ReplyTime = config.ConfigStruct.Xhh.ReplyTime
	if CheckTime == 0 {
		loger.Loger.Warn("[XHH]您的设置中未设置检查时间，已默认为30s")
		CheckTime = 30
	}
	if ReplyTime == 0 {
		loger.Loger.Warn("[XHH]您的设置中未设置回复间隔，已默认为10s")
		ReplyTime = 10
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

var DontReply bool

func CheckAt() {
	fmt.Println("[XHH]检查@", time.Now().Format("2006-01-02 15:04:05"))
	var offset int
	nomore := "false"
	other := fmt.Sprintf("?message_type=16&offset=%v&limit=20&no_more=%s", offset, nomore)
	resp := SendReq("GET", "/bbs/app/user/message", nil, other)
	if resp == nil {
		loger.Loger.Error("[XHH]链接发送失败了")
		return
	}
	var data Respo
	Dbyte, err := io.ReadAll(resp.Body)
	if err != nil {
		loger.Loger.Error("[XHH]无法读取Body", zap.Error(err))
		return
	}
	err = json.Unmarshal(Dbyte, &data)

	if err != nil {
		loger.Loger.Error("[XHH]无法反序列化", zap.Error(err), zap.String("raw", string(Dbyte)))
		return
	}

	for _, v := range data.Result.Messages {
		if Check(v.UserID) {
			if DontReply {
				db.Insert(v.MsgID, v.CommentID, v.RootCommentID, v.LinkID, v.UserID, v.CommentText, true)
			} else {
				db.Insert(v.MsgID, v.CommentID, v.RootCommentID, v.LinkID, v.UserID, v.CommentText, false)
			}
		}
	}
	DontReply = false
	time.Sleep(time.Duration(CheckTime) * time.Second)
	CheckAt()
}

func AutoReply() {
	for {
		Arr := db.GetComm()
		if len(Arr) == 0 {
			fmt.Println("[XHH]无可回复", time.Now().Format("2006-01-02 15:04:05"))
			time.Sleep(time.Duration(ReplyTime) * time.Second)
			continue
		}
		var wg sync.WaitGroup
		loger.Loger.Info("[XHH]正在回复评论", zap.Int("评论数", len(Arr)))
		wg.Add(len(Arr))
		for _, v := range Arr {
			go func() {
				defer wg.Done()
				if v.CommentID != 0 {
					var isok bool
					if Check(v.Uid) {
						Info, top, tags := GetLinkInfo(v.LinkID)
						if len(Info) <= 1 {
							loger.Loger.Info("[XHH]获取LinkID失败")
							return
						}
						ReplyText := ai.GetAiReply(Info, v.Text, top, tags)
						if ReplyText == "" {
							loger.Loger.Info("[XHH]Ai返回错误")
							return
						}
						isok = Reply(ReplyText, strconv.Itoa(v.LinkID), strconv.Itoa(v.CommentID), strconv.Itoa(v.RootID), "")

					}
					if isok {
						db.Replyed(v.CommentID)
					} else {
						loger.Loger.Error("[XHH]无法回复评论")
					}
				} else {
					wg.Done()
					fmt.Println("[XHH]无事可做")
				}
			}()
		}
		wg.Wait()
	}

}
