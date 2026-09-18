package log

import (
	"log"
	"os"
)

type writer struct{}

var Loger *log.Logger
var Mod string

func (w writer) Write(b []byte) (n int, err error) {
	go os.Stdout.Write(b)
	n = len(b)
	return
}

func Init() {
	if Mod == "debug" {
		DebugerInit()
	} else {
		var w writer
		Loger = log.New(w, "[log]", log.Ldate|log.Ltime|log.Lmsgprefix)
		Loger.Printf("日志初始化完毕")
	}
}

func DebugerInit() {
	var w writer
	Loger = log.New(w, "[log|debug]", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
	Loger.Printf("日志初始化完毕")
}

func DebugLog(arg ...any) {
	go func() {
		if Mod != "debug" {
			return
		}
		Loger.Println(arg...)
	}()
}
