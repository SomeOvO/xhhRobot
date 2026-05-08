package config

import (
	"encoding/json"
	"os"
	"xhhrobot/loger"
)

var ConfigStruct struct {
	Xhh struct {
		Owner   int    `json:"owner"`
		BaseUrl string `json:"baseUrl"`
		WebVer  string `json:"webver"`
		Ver     string `json:"version"`
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
			loger.Loger.Fatal("Plz edit config and restart")
		}
		panic(err)
	}
	err = json.Unmarshal(file, &ConfigStruct)
	if err != nil {
		panic(err)
	}
	loger.Loger.Info("[CFG]Init OK")
}
