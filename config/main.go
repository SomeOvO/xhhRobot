package config

import (
	"encoding/json"
	"os"
	"xhhrobot/loger"
)

var ConfigStruct struct {
	Xhh struct {
		CheckTime int    `json:"checkTime"`
		ReplyTime int    `json:"replyTime"`
		Owner     string `json:"owner"`
		DeviceID  string `json:"deviceID"`
		BaseUrl   string `json:"baseUrl"`
		WebVer    string `json:"webver"`
		Ver       string `json:"version"`
	} `json:"xhh"`
	DataBase struct {
		Type   string `json:"type"`
		Db     string `json:"db"`
		Host   string `json:"host"`
		Port   string `json:"port"`
		User   string `json:"user"`
		Passwd string `json:"passwd"`
	} `json:"database"`
	Ai struct {
		Model   string `json:"model"`
		Prompt  string `json:"prompt"`
		BaseUrl string `json:"baseUrl"`
		Token   string `json:"token"`
	} `json:"ai"`
}

func InitConfig() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	file, err := os.ReadFile(wd + "/config.json")
	if err != nil {
		if os.IsNotExist(err) {
			Data, err := json.Marshal(ConfigStruct)
			if err != nil {
				panic(err)
			}
			os.WriteFile("./config.json", Data, 0644)
			loger.Loger.Fatal("请修改配置文件后重新启动")
		}
		panic(err)
	}
	err = json.Unmarshal(file, &ConfigStruct)
	if err != nil {
		panic(err)
	}
	loger.Loger.Info("[CFG]Init OK")
}
