package xhh

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
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
		loger.Loger.Fatal("[XHH]未检测到Cookie，请先登陆")
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
var errInfo struct {
	Count   int
	LastErr time.Time
}

func IsErr() {
	now := time.Now()

	if now.Sub(errInfo.LastErr) > 10*time.Minute {
		errInfo.Count = 1
		errInfo.LastErr = now
		return
	}

	errInfo.Count++

	if errInfo.Count >= 5 {
		loger.Loger.Fatal("[XHH]程序十分钟内错误五次，已退出防止频繁错误")
	}
}

func CheckAt() {
	// 依赖api返回的msg_id是从新到旧，即msg_id从大到小
	lastMsgID := 0 // 储存上一轮查询的最大msg_id，也就是上一轮查询过的最新的

	// 初始化lastMsgID
	if db.IsNew() {
		if DontReply {
			other := fmt.Sprintf("?message_type=16&offset=%v&limit=%v&no_more=%s", 0, 1, "false")
			resp := SendReq("GET", "/bbs/app/user/message", nil, other)
			if resp == nil {
				loger.Loger.Error("[XHH]链接发送失败了")
				IsErr()
			}

			var data Respo
			Dbyte, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				loger.Loger.Error("[XHH]无法读取Body", zap.Error(err))
				IsErr()
			}

			err = json.Unmarshal(Dbyte, &data)
			if err != nil {
				loger.Loger.Error("[XHH]无法反序列化", zap.Error(err), zap.String("raw", string(Dbyte)))
				IsErr()
			}

			// 若没有@，则初始化为0
			// 上面出错不退出，也初始化为0，就是会翻页翻到底
			if len(data.Result.Messages) > 0 {
				lastMsgID = data.Result.Messages[0].MsgID
			}
		}

		// 需要回复的情况由于lastMsgID初始化为零，会通过不断翻页处理
	} else {
		lastMsgID = db.GetNewestMsgIDNotReplied()
	}
	fmt.Println("[XHH]初始化lastMsgID为", lastMsgID, time.Now().Format("2006-01-02 15:04:05"))

	for {
		fmt.Println("[XHH]检查@", time.Now().Format("2006-01-02 15:04:05"))
		limit := 20
		nomore := "false"
		newLastMsgID := lastMsgID // 用于储存这一次查询的最新msg_id
		curMsgID := math.MaxInt   // 用于遍历这一轮的msg，直到上一轮的查询的最新msg_id为止

		for offset := 0; curMsgID > lastMsgID; offset += limit {
			other := fmt.Sprintf("?message_type=16&offset=%v&limit=%v&no_more=%s", offset, limit, nomore)
			resp := SendReq("GET", "/bbs/app/user/message", nil, other)
			if resp == nil {
				loger.Loger.Error("[XHH]链接发送失败了")
				IsErr()
				continue
			}

			var data Respo
			Dbyte, err := io.ReadAll(resp.Body)
			resp.Body.Close()

			if err != nil {
				loger.Loger.Error("[XHH]无法读取Body", zap.Error(err))
				IsErr()
				continue
			}
			err = json.Unmarshal(Dbyte, &data)
			if err != nil {
				loger.Loger.Error("[XHH]无法反序列化", zap.Error(err), zap.String("raw", string(Dbyte)))
				IsErr()
				continue
			}

			msgsLen := len(data.Result.Messages)
			if msgsLen == 0 {
				// 处理lastMsgID==0的情况，遍历页直到api返回空退出
				break
			}
			if offset == 0 {
				newLastMsgID = data.Result.Messages[0].MsgID
			}
			curMsgID = data.Result.Messages[msgsLen-1].MsgID // 更新为这一页最小的msg_id

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
			// time.Sleep(time.Duration(3) * time.Second)
			// 不在此处sleep应该不会风控
			// limit设置大一些大可改善，但是感觉应该没那么多人用
		}

		lastMsgID = newLastMsgID
		time.Sleep(time.Duration(CheckTime) * time.Second)
	}
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
						Info, top, tags := GetLinkInfo(v.LinkID, v.CommentID)
						if len(Info) <= 1 {
							loger.Loger.Info("[XHH]获取LinkID失败")
							IsErr()
							return
						}
						ReplyText := ai.GetAiReply(Info, v.Text, top, tags)
						if ReplyText == "" {
							loger.Loger.Info("[XHH]Ai返回错误")
							IsErr()
							return
						}
						isok = Reply(ReplyText, strconv.Itoa(v.LinkID), strconv.Itoa(v.CommentID), strconv.Itoa(v.RootID), "")

					}
					if isok {
						db.Replyed(v.CommentID)
					} else {
						IsErr()
						loger.Loger.Error("[XHH]无法回复评论")
					}
				} else {
					fmt.Println("[XHH]无事可做")
				}
			}()
		}
		wg.Wait()
		time.Sleep(time.Duration(ReplyTime) * time.Second)
	}
}
