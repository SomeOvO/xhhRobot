package sqlite

import (
	"database/sql"
	"heybox/database"
	"heybox/log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func Init() {
	log.DebugLog("数据库为：SQLite")
	database.Init = connect
}

func connect() {
	var err error
	db, err = sql.Open("sqlite", "./sql.db")
	if err != nil {
		log.Loger.Fatal("无法链接至数据", err)
	}
	err = db.Ping()
	if err != nil {
		log.Loger.Fatal("测试数据库链接失败", err)
	}
}
