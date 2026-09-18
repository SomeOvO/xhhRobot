package main

import (
	"heybox/cfg"
	xhh_methods "heybox/xhh/methods"
)

func main() {
	cfg.Config.Xhh.BaseUrl = "https://api.xiaoheihe.cn"
	cfg.Config.Xhh.Version = "999.0.4"
	cfg.Config.Xhh.WebVer = "2.5"
	xhh_methods.Login()
}
