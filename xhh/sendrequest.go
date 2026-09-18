package xhh

import (
	"errors"
	"fmt"
	"heybox/cfg"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type RequstBody struct {
	Method *string
	Path   string
	Url    *url.URL
	Body   io.Reader
	Query  *string
}

const (
	GET string = "GET"
)

func SendRequst(RequstBody *RequstBody) (*http.Response, *[]byte, error) {
	if err := makeBody(RequstBody); err != nil {
		return nil, nil, err
	}
	ReqData := *RequstBody
	request, err := http.NewRequest(*ReqData.Method, ReqData.Url.String(), ReqData.Body)
	if err != nil {
		return nil, nil, err
	}
	request.Header.Set("host", ReqData.Url.Host)
	request.Header.Set("Referer", fmt.Sprintf("%s://%s", ReqData.Url.Scheme, ReqData.Url.Host))

	if ReqData.Body != nil {
		request.Header.Set("content-type", "application/x-www-form-urlencoded;charset=utf-8")
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	Data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return resp, &Data, nil
}

// 构建请求体
func makeBody(Body *RequstBody) error {
	XhhConfig := cfg.Config.Xhh
	DefaultMethod := "GET"
	DefaultUrl := XhhConfig.BaseUrl
	if Body.Method == nil {
		Body.Method = &DefaultMethod
	}

	if Body.Path == "" {
		return errors.New("请求路径不可为空")
	}
	if Body.Url == nil {
		RequstUrl, err := url.Parse(DefaultUrl)
		if err != nil {
			return err
		}
		RequstUrl.Path = Body.Path
		Body.Url = RequstUrl
	}
	if Body.Query != nil {
		Body.Url.RawQuery = *Body.Query
	}
	hkey, nonce, time, err := GetKeys(Body.Path)
	if err != nil {
		return err
	}
	query := Body.Url.Query()
	query.Set("os_type", "web")
	query.Set("app", "web")
	query.Set("client_type", "web")
	query.Set("version", XhhConfig.Version)
	query.Set("web_version", XhhConfig.WebVer)
	query.Set("x_client_type", "web")
	query.Set("x_app", "heybox_website")
	query.Set("x_os_type", "Windows")
	query.Set("device_info", "Chrome")
	query.Set("device_id", XhhConfig.DeviceID)
	query.Set("hkey", hkey)
	query.Set("_time", strconv.Itoa(time))
	query.Set("nonce", nonce)
	query.Set("_notip", "true")
	Body.Url.RawQuery = query.Encode()
	return nil
}
