package main

import (
	"heybox/cfg"
	"heybox/database"
	"heybox/database/sqlite"
	"heybox/log"
)

func main() {
	cfg.Config.Xhh.BaseUrl = "https://api.xiaoheihe.cn"
	cfg.Config.Xhh.Version = "999.0.4"
	cfg.Config.Xhh.WebVer = "2.5"
	log.Mod = "debug"
	log.Init()
	sqlite.Init()
	database.Init()
}
