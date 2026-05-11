package xhh

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"xhhrobot/config"
	"xhhrobot/loger"

	"go.uber.org/zap"
)

func SendReq(Method, Path string, Body io.Reader, other string) *http.Response {
	cfg := config.ConfigStruct.Xhh
	u, err := url.Parse(cfg.BaseUrl + Path + other)
	if err != nil {
		loger.Loger.Error("[SendReq]Creat requset url failed")
		return nil
	}
	reqUrl := u.Query()
	hkey, nonce, time := GetKeys(Path)
	reqUrl.Set("os_type", "web")
	reqUrl.Set("app", "web")
	reqUrl.Set("client_type", "web")
	reqUrl.Set("version", cfg.Ver)
	reqUrl.Set("web_version", cfg.WebVer)
	reqUrl.Set("x_client_type", "web")
	reqUrl.Set("x_app", "heybox_website")
	if Info.HeyBoxId != "" {
		reqUrl.Set("heybox_id", Info.HeyBoxId)
	}
	reqUrl.Set("x_os_type", "Windows")
	reqUrl.Set("device_info", "Chrome")
	reqUrl.Set("device_id", cfg.DeviceID)
	reqUrl.Set("hkey", hkey)
	reqUrl.Set("_time", strconv.Itoa(time))
	reqUrl.Set("nonce", nonce)
	reqUrl.Set("_notip", "true")
	u.RawQuery = reqUrl.Encode()
	req, err := http.NewRequest(Method, u.String(), Body)
	if err != nil {
		loger.Loger.Error("[SendReq]Can't Create HttpReq", zap.Error(err))
		return nil
	}
	req.Header.Set("host", "api.xiaoheihe.cn")
	if Info.Cookie != "" {
		req.Header.Set("cookie", Info.Cookie)
	}
	req.Header.Set("Referer", "https://www.xiaoheihe.cn/")
	if Body != nil {
		req.Header.Set("content-type", "application/x-www-form-urlencoded;charset=utf-8")
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		loger.Loger.Error("[XHH] SendReq Failed", zap.Error(err))
	}
	return resp
}
