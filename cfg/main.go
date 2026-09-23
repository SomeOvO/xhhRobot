package cfg

import "heybox/database"

type config map[int]*database.TYPE_Config

var Config config

//XhhConfig 0-100
const (
	CONFIG_XHH_REPLYURL = 0
	CONFIG_XHH_BASEURL  = 1
	CONFIG_XHH_VERSION  = 2
	CONFIG_XHH_WEBVER   = 3
	CONFIG_XHH_DEVICEID = 4
)

//AiConfig 101-200
const (
	CONFIG_AI_MODEL   = 101
	CONFIG_AI_BaseUrl = 102
	CONFIG_AI_TOKEN   = 103
)

func SetConfig(configInfo *database.TYPE_Config) error {
	var Set = database.Db.Config.Set
	Config[configInfo.Name] = configInfo
	return Set(configInfo.Name, configInfo.Value, configInfo.Enable)

}

// 加载配置文件
func LoadConfig() {
	Config = make(config)
	//加载小黑盒配置
}
func loadxhhconfig() {
	var Get = database.Db.Config.Get

	Config[CONFIG_XHH_REPLYURL] = Get(CONFIG_XHH_REPLYURL)
	Config[CONFIG_XHH_BASEURL] = Get(CONFIG_XHH_BASEURL)
	Config[CONFIG_XHH_VERSION] = Get(CONFIG_XHH_VERSION)
	Config[CONFIG_XHH_WEBVER] = Get(CONFIG_XHH_WEBVER)
	Config[CONFIG_XHH_DEVICEID] = Get(CONFIG_XHH_DEVICEID)

	//设置默认
	if Config[CONFIG_XHH_REPLYURL] == nil {

	}
}
