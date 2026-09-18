package xhh_methods

import (
	"encoding/json"
	"errors"
	"fmt"
	"heybox/xhh"
	"net/http"
	"net/url"
	"time"

	"github.com/skip2/go-qrcode"
)

type login_body struct {
	Status  string `json:"status"`
	Msg     string `json:"msg"`
	Version string `json:"version"`
	Result  struct {
		Qrcode   string `json:"qr_url"`
		Expire   int    `json:"expire"`
		ErrMsg   string `json:"error_msg"`
		Err      string `json:"error"`
		NickName string `json:"nickname"`
	} `json:"result"`
}

func Login() error {
	//获取qrcode的URL
	Qrcode_url, err := getqrcode()
	if err != nil {
		return err
	}
	url, err := url.Parse(Qrcode_url)
	if err != nil {
		return err
	}
	qr := url.Query().Get("qr")
	if qr == "" {
		return errors.New("获取到的qr为空")
	}
	ticker := time.Tick(1 * time.Second)
	start := time.Now().Unix()
	for range ticker {
		if time.Now().Unix()-start > 300 {
			return errors.New("登陆超时")
		}
		err, isok := check(qr)
		if err != nil {
			return err
		}
		if isok {
			break
		}
	}
	return nil
}

func getqrcode() (url string, err error) {
	var RequstBody xhh.RequstBody
	RequstBody.Path = "/account/get_qrcode_url/"
	_, Resp, err := xhh.SendRequst(&RequstBody)
	if err != nil {
		return
	}
	var Body login_body
	if err = json.Unmarshal(*Resp, &Body); err != nil {
		return
	}
	if Body.Status != "ok" {
		return "", fmt.Errorf("请求登陆失败,%v", string(*Resp))
	}
	code, err := qrcode.New(Body.Result.Qrcode, qrcode.High)
	if err != nil {
		return
	}
	err = code.WriteFile(256, "qrcode.png")
	if err != nil {
		return
	}
	fmt.Println(code.ToSmallString(true))
	fmt.Println("如无法扫码，可查看运行目录下的qrcode.png")
	return Body.Result.Qrcode, nil
}
func check(qr string) (err error, isok bool) {
	var RequstBody xhh.RequstBody
	RequstBody.Path = "/account/qr_state/"
	query := "qr=" + qr
	RequstBody.Query = &query
	_, respBytes, err := xhh.SendRequst(&RequstBody)
	if err != nil {
		return
	}
	var body login_body
	if err = json.Unmarshal(*respBytes, &body); err != nil {
		return
	}

	if body.Result.Err == "ok" {
		fmt.Println("\r\n登陆成功")
		fmt.Printf("欢迎 > %s", body.Result.NickName)
		isok = true
	} else {
		msg := body.Result.ErrMsg
		if msg == "" {
			msg = "等待扫码"
		}
		fmt.Printf("\r%s %v", msg, time.Now().Format("2006-01-02 15:04:05"))
	}
	return
}
func savecookie(Cookies []*http.Cookie) {

}
