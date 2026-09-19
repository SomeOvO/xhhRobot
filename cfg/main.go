package cfg

var Config configstruct

type configstruct struct {
	Xhh xhhconfig
}
type config struct {
	Name   string
	Value  string
	Enable bool
}
type xhhconfig map[int]config

const (
	CONFIG_XHH_BASEURL  = 1
	CONFIG_XHH_VERSION  = 2
	CONFIG_XHH_WEBVER   = 3
	CONFIG_XHH_DEVICEID = 4
)

// 加载配置文件
func LoadConfig() {

}
