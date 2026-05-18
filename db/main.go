package db

import (
	"context"
	"xhhrobot/config"
	"xhhrobot/loger"
	"xhhrobot/pg"
	"xhhrobot/sqlite"

	"go.uber.org/zap"
)

var cfg = &config.ConfigStruct.DataBase

func Init() {
	switch cfg.Type {
	case "pg":
		pg.InitPostgreSQL()
		return
	case "sqlite":
		sqlite.Init()
		return
	default:
		loger.Loger.Fatal("[DB]无效的数据库类型")
	}
}

func Insert(msg_id, comment_a_id, comment_root_id, link_id, user_a_id int, comment_text string, reply bool) bool {
	ctx := context.Background()
	if cfg.Type == "pg" {
		_, err := pg.Conn.Exec(ctx, "INSERT INTO at (msg_id,comment_a_id,comment_root_id,link_id,user_a_id,comment_text,reply) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (msg_id) DO NOTHING", msg_id, comment_a_id, comment_root_id, link_id, user_a_id, comment_text, reply)
		if err != nil {
			loger.Loger.Info("[DB]PsqlError", zap.Error(err))
			return false
		}
		return true
	}
	if cfg.Type == "sqlite" {
		_, err := sqlite.Db.Exec("INSERT INTO at (msg_id,comment_a_id,comment_root_id,link_id,user_a_id,comment_text,reply) VALUES (?,?,?,?,?,?,?) ON CONFLICT (msg_id) DO NOTHING", msg_id, comment_a_id, comment_root_id, link_id, user_a_id, comment_text, reply)
		if err != nil {
			loger.Loger.Info("[DB]SQLiteERROR", zap.Error(err))
			return false
		}
	}
	return false
}

func Replyed(comment_id int) {
	ctx := context.Background()
	if cfg.Type == "pg" {
		pg.Conn.Exec(ctx, "UPDATE at SET reply=$1 WHERE comment_a_id=$2", true, comment_id)
	}
	if cfg.Type == "sqlite" {
		sqlite.Db.Exec("UPDATE at SET reply=? WHERE comment_a_id=?", true, comment_id)
	}
}

type CommStruct struct {
	LinkID    int
	CommentID int
	RootID    int
	Text      string
	Uid       int
}

func GetComm() (CommArr []CommStruct) {
	ctx := context.Background()
	if cfg.Type == "pg" {
		row, err := pg.Conn.Query(ctx, "SELECT link_id,comment_a_id,comment_root_id,comment_text,user_a_id FROM at WHERE reply=false AND comment_a_id<>0 AND link_id<>0 LIMIT 3")
		if err != nil {
			loger.Loger.Error("[DB]无法获取评论信息", zap.Error(err))
			return
		}
		for row.Next() {
			var Comm CommStruct
			row.Scan(&Comm.LinkID, &Comm.CommentID, &Comm.RootID, &Comm.Text, &Comm.Uid)
			CommArr = append(CommArr, Comm)
		}
		return
	}
	if cfg.Type == "sqlite" {
		row, err := sqlite.Db.Query("SELECT link_id,comment_a_id,comment_root_id,comment_text,user_a_id FROM at WHERE reply=false AND comment_a_id<>0 AND link_id<>0 LIMIT 3")
		if err != nil {
			loger.Loger.Error("[DB]无法获取评论信息", zap.Error(err))
			return
		}
		for row.Next() {
			var Comm CommStruct
			row.Scan(&Comm.LinkID, &Comm.CommentID, &Comm.RootID, &Comm.Text, &Comm.Uid)
			CommArr = append(CommArr, Comm)
		}
	}

	return
}

func IsNew() bool {
	ctx := context.Background()
	var num int
	if cfg.Type == "pg" {
		row := pg.Conn.QueryRow(ctx, "SELECT COUNT(*) FROM at")
		row.Scan(&num)
	}
	if cfg.Type == "sqlite" {
		row := sqlite.Db.QueryRow("SELECT COUNT(*) FROM at")
		row.Scan(&num)
	}
	if num > 0 {
		return false
	} else {
		return true
	}
}
