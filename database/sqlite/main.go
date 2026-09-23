package sqlite

import (
	"database/sql"
	"heybox/database"
	"heybox/log"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB
var sqlite database.DataBase

func Init() {
	log.DebugLog("数据库为：SQLite")
	sqlite.Conn = connect

	var sqlite_config database.Config
	sqlite_config.Set = config_set
	sqlite_config.Get = config_get
	sqlite.Config = sqlite_config

	var sqlite_xhh database.Xhh
	sqlite_xhh.GetCookie = xhh_getcookie
	sqlite_xhh.SaveCookie = xhh_savecookie
	sqlite.Xhh = sqlite_xhh
	sqlite.Remove = remove
	database.Db = &sqlite
}

func connect() error {
	var err error
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	db, err = sql.Open("sqlite", wd+"/sql.db")
	if err != nil {
		log.Loger.Fatal("无法链接至数据", err)
		return err
	}
	err = db.Ping()
	if err != nil {
		log.Loger.Fatal("测试数据库链接失败", err)
		return err
	}
	return createDatabase()
}

func createDatabase() error {
	for _, v := range database.DatabaseList {
		_, err := db.Exec(v)
		if err != nil {
			return err
		}
	}
	return nil
}
