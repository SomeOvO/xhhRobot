package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
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

func CheckNew() {
	if !db.IsNew() {
		return
	}
	fmt.Println("检测到您是第一次运行\n是否允许将先前@过的名单加入至艾特列表？\ny(es) or n(o) 默认n\n请输入y或n")
	input := bufio.NewReader(os.Stdin)
	str, err := input.ReadString('\n')
	if err != nil {
		loger.Loger.Fatal("[MAIN]无法读取您的输入")
	}
	in := strings.TrimRight(str, "\r\n")

	switch in {
	case "y":
		xhh.DontReply = false
		return
	case "yes":
		xhh.DontReply = false
		return
	case "Y":
		xhh.DontReply = false
		return
	case "Yes":
		xhh.DontReply = false
		return
	case "YES":
		xhh.DontReply = false
		return
	default:
		xhh.DontReply = true
		return
	}
}

func start(mode *string) {
	switch *mode {
	case "default":
		loger.Loger.Info("\nhttps://github.com/SomeOvO/xhhRobot\n你需要输入启动项\n-mode start | login | test")
	case "test":
		xhh.RunTest()
	case "login":
		xhh.Login()
	case "start":
		CheckNew()
		xhh.Init()
		xhh.Start()
		select {}
	}

}

