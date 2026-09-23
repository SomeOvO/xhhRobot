package database

//数据库，新增的数据库都要支持以下操作

type DataBase struct {
	Conn   func() error
	Config Config
	Xhh    Xhh
	//Remove函数不应在测试外被调用，如需调用请输入对应的key
	Remove func(key string) error
}
type Xhh struct {
	SaveCookie func(uid int, pkey, token string) error
	GetCookie  func(uid int) (err error, pkey, token string)
}
type Config struct {
	Set func(Name int, Value string, Enable bool) error
	Get func(ConfigName int) (config *TYPE_Config)
}

var Db *DataBase
