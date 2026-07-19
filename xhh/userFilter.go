package xhh

import (
	"strconv"
	"strings"
	"sync"

	"xhhrobot/config"
	"xhhrobot/loger"
)

var (
	userFilterList map[int]struct{}
	initUserFilter = sync.OnceFunc(func() {
		cfg := config.ConfigStruct.Xhh.UserFilter
		userFilterList = make(map[int]struct{})
		uidStrs := strings.Split(cfg.UserIDs, ",")
		for _, uidStr := range uidStrs {
			if uidStr != "" {
				i, err := strconv.Atoi(uidStr)
				if err != nil {
					loger.Loger.Error("[XHH]您的过滤名单配置->" + uidStr + "<-似乎并非数字")
					continue
				}
				userFilterList[i] = struct{}{}
			}
		}
		switch cfg.Mode {
		case "blacklist", "whitelist":
		default:
			loger.Loger.Fatal("[XHH]过滤名单模式应为blacklist或whitelist")
		}
		if cfg.Mode == "whitelist" && len(userFilterList) == 0 {
			loger.Loger.Warn("[XHH]您的过滤名单模式为白名单，但名单为空")
		}
	})
)

func Check(UID int) bool {
	initUserFilter()
	_, exists := userFilterList[UID]
	if config.ConfigStruct.Xhh.UserFilter.Mode == "blacklist" {
		return !exists
	}
	return exists
}
