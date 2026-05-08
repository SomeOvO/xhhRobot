package main

import (
	"flag"
	"time"
	"xhhrobot/config"
	"xhhrobot/db"
	"xhhrobot/loger"
	"xhhrobot/xhh"
)

func main() {
	loger.InitLog()
	config.InitConfig()
	time.Sleep(1 * time.Second)
	db.Init()
	mode := flag.String("mode", "default", "Switch a mode when start")
	flag.Parse()
	start(mode)
}

func start(mode *string) {
	switch *mode {
	case "default":
		loger.Loger.Info("\nHi,This is XhhRobot\nYou need start with a mode\n-mode start | login | test")
	case "test":
		xhh.GetLinkInfo(180970107)
	case "login":
		xhh.Login()
	case "start":
		xhh.Init()
		xhh.Start()
		select {}
	}

}
