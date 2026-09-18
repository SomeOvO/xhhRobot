package cfg

type configstruct struct {
	Xhh xhhconfig
}

type xhhconfig struct {
	BaseUrl  string
	Version  string
	WebVer   string
	DeviceID string
}

var Config configstruct
