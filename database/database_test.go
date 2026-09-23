package database_test

import (
	"heybox/database"
	"heybox/database/sqlite"
	"heybox/log"
	"testing"
)

func TestMain(T *testing.M) {
	log.Init()
	sqlite.Init()
	T.Run()
}
func TestSqliteConn(T *testing.T) {
	db := database.Db
	if db == nil {
		T.Fatalf("数据库为空")
	}
	T.Log("创建数据库链接中", db)
	if err := db.Conn(); err != nil {
		T.Fatal("数据库链接错误", err)
	}
	if err := db.Config.Set(1, "2", true); err != nil {
		T.Fatal("测试更新配置失败", err)
	}
	res := db.Config.Get(1)
	T.Log(res, res.Name, res.Value, res.Enable)
	if err := db.Remove("Dont Input THIS"); err != nil {
		T.Fatal("移除数据库出错", err)
	}
	T.Log("SQLite测试成功")
}

func TestSqliteXHH(t *testing.T) {
	db := database.Db
	err := db.Conn()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("存入数据库")
	err = db.Xhh.SaveCookie(123, "456", "789")
	if err != nil {
		t.Fatal(err)
	}
	err, pkey, token := db.Xhh.GetCookie(123)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pkey, token)
	if pkey != "456" || token != "789" {
		t.Fatal("读取失败")
	}
	t.Log("测试数据库冲突")
	err = db.Xhh.SaveCookie(123, "123", "456")
	if err != nil {
		t.Fatal(err)
	}
	err, pkey, token = db.Xhh.GetCookie(123)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pkey, token)
	if pkey != "123" || token != "456" {
		t.Fatal("冲突测试失败")
	}
}
