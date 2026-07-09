package sqlite

import (
	"database/sql"
	"xhhrobot/loger"

	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

var Db *sql.DB

func Init() {
	var err error
	Db, err = sql.Open("sqlite", "./sql.db")
	if err != nil {
		loger.Loger.Fatal("[SQLite]无法读取文件", zap.Error(err))
	}
	_, err = Db.Exec(`
	CREATE TABLE IF NOT EXISTS at (
	msg_id BIGINT PRIMARY KEY,
	comment_a_id BIGINT,
	comment_root_id BIGINT,
	link_id BIGINT,
	user_a_id BIGINT,
	comment_text TEXT,
	reply boolean
	)
	`)
	if err != nil {
		loger.Loger.Fatal("[Sqlite]无法创建新的数据库", zap.Error(err))
	}
	err = Db.Ping()
	if err != nil {
		loger.Loger.Fatal("[Sqlite]无法连接至新的数据库", zap.Error(err))
	}
	loger.Loger.Info("[SQLite]READY!")
}
