package database_test

import (
	"heybox/database"
	"heybox/database/sqlite"
	"testing"
)

func TestSqlite(T *testing.T) {
	sqlite.Init()
	db := database.Db
	if db == nil {
		T.Fatalf("数据库为空")
	}
	T.Log("创建数据库链接中", db)
	if err := db.Conn(); err != nil {
		T.Fatal("数据库链接错误", err)
	}
	if err := db.Remove("Dont Input THIS"); err != nil {
		T.Fatal("移除数据库出错", err)
	}
	T.Log("SQLite测试成功")
}
