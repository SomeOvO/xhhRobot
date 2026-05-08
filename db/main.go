package db

import (
	"context"
	"xhhrobot/config"
	"xhhrobot/loger"
	"xhhrobot/pg"

	"go.uber.org/zap"
)

var cfg = &config.ConfigStruct.DataBase

func Init() {
	switch cfg.Type {
	case "pg":
		pg.InitPostgreSQL()
		return
	case "sqlite":
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
	return false
}

func Replyed(comment_id int) {
	ctx := context.Background()
	if cfg.Type == "pg" {
		pg.Conn.Exec(ctx, "UPDATE at SET reply=$1 WHERE comment_a_id=$2", true, comment_id)
	}
}

func GetComm() (linkID, commentID, rootID int, text string, UID int) {
	ctx := context.Background()
	if cfg.Type == "pg" {
		row := pg.Conn.QueryRow(ctx, "SELECT link_id,comment_a_id,comment_root_id,comment_text,user_a_id FROM at WHERE reply=false LIMIT 1")
		row.Scan(&linkID, &commentID, &rootID, &text, &UID)
		return
	}
	return
}
