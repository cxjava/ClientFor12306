package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Login struct {
	Username string
	Password string
	Captcha  string
	Cookie   string
	CDN      string
}

type CheckRandCode struct {
	Basic
	Data             string      `json:"data"`
	Messages         []string    `json:"messages,omitempty"`
	ValidateMessages interface{} `json:"validateMessages,omitempty"`
}

type LoginAysnSuggest struct {
	Basic
	Data struct {
		LoginCheck string `json:"loginCheck"`
	}
	Messages         []string    `json:"messages,omitempty"`
	ValidateMessages interface{} `json:"validateMessages,omitempty"`
}

//获取新的cookie
func (l *Login) setNewCookie() error {
	if l.Cookie != "" {
		return nil
	}
	req, err := http.NewRequest("GET", URLLoginJs, nil)
	if err != nil {
		Error("setNewCookie http.NewRequest error:", err)
		return err
	}

	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/login/init"}
	AddReqestHeader(req, "GET", h)

	resp, err := client.Do(req)
	if err != nil {
		Error("setNewCookie client.Do error:", err)
		return err
	}
	defer resp.Body.Close()

	l.Cookie = GetCookieFromRespHeader(resp)
	Info("Cookie=" + l.Cookie + "=")
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Error("ioutil.ReadAll:", err)
		return err
	}
	Info(string(bodyByte))
	return nil
}

func (l *Login) CheckRandCodeAnsyn() (r bool, msg []string) {
	b := url.Values{}
	b.Add("randCode", l.Captcha)
	b.Add("rand", Rand)
	params, err := url.QueryUnescape(b.Encode())
	if err != nil {
		Error("CheckRandCodeAnsyn url.QueryUnescape error:", err)
		return false, []string{err.Error()}
	}
	Info("CheckRandCodeAnsyn params:", params)
	content, err := DoForWardRequest(login.CDN, "POST", URLCheckRandCodeAnsyn, strings.NewReader(params))
	if err != nil {
		Error("CheckRandCodeAnsyn DoForWardRequest error:", err)
		return false, []string{err.Error()}
	}
	crc := &CheckRandCode{}

	if err := json.Unmarshal([]byte(content), &crc); err != nil {
		Error("CheckRandCodeAnsyn json.Unmarshal error:", err)
		return false, []string{err.Error()}
	}
	Info(crc)
	return crc.Data == "Y", crc.Messages
}

func (l *Login) checkUser() (bool, error) {
	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/leftTicket/init"}
	body, err := DoForWardRequestHeader(l.CDN, "POST", URLCheckUser, nil, h)
	if err != nil {
		Error("checkUser DoForWardRequest error:", err)
		return false, err
	}
	Info("checkUser body:", body)
	return strings.Contains(body, `"flag":true`), nil
}

func (l *Login) loginAysnSuggest() (r bool, msg []string) {
	b := url.Values{}
	b.Add("loginUserDTO.user_name", l.Username)
	b.Add("userDTO.password", l.Password)
	b.Add("randCode", l.Captcha)
	params, err := url.QueryUnescape(b.Encode())
	if err != nil {
		Error("Login url.QueryUnescape error:", err)
		return false, []string{err.Error()}
	}
	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/login/init"}
	content, err := DoForWardRequestHeader(l.CDN, "POST", URLLoginAysnSuggest, strings.NewReader(params), h)
	if err != nil {
		Error("CheckRandCodeAnsyn DoForWardRequest error:", err)
		return false, []string{err.Error()}
	}

	las := &LoginAysnSuggest{}
	if err := json.Unmarshal([]byte(content), &las); err != nil {
		Error("Login json.Unmarshal error:", err)
		return false, []string{err.Error()}
	}
	Info(las)
	return las.Data.LoginCheck == "Y", las.Messages
}

func (l *Login) userLogin() {
	val := url.Values{}
	val.Add("_json_att", Json_att)

	params, _ := url.QueryUnescape(val.Encode())
	Info("userLogin params:", params)
	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/login/init"}
	body, err := DoForWardRequestHeader(l.CDN, "POST", URLUserLogin, strings.NewReader(params), h)
	if err != nil {
		Error("userLogin DoForWardRequest error:", err)
	}
	Debug("userLogin body:", body)
}

//获取联系人
func (l *Login) getPassengerDTO() (p PassengerDTO) {
	val := url.Values{}
	params, _ := url.QueryUnescape(val.Encode())
	Info("getPassengerDTO params:", params)
	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/leftTicket/init"}
	body, err := DoForWardRequestHeader(l.CDN, "POST", URLGetPassengerDTOs, strings.NewReader(params), h)
	if err != nil {
		Error("getPassengerDTO DoForWardRequest error:", err)
		return
	}
	Debug("getPassengerDTO body:", body)
	if !strings.Contains(body, "passenger_name") {
		Error("获取联系人出错!!!!!!返回:", body)
		return
	}
	p = PassengerDTO{}
	if err := json.Unmarshal([]byte(body), &p); err != nil {
		Error("getPassengerDTO", l.CDN, err)
		return
	}
	return
}

func (l *Login) initQueryUserInfo() {
	time.Sleep(time.Second * 5)
	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/leftTicket/init"}
	body, err := DoForWardRequestHeader(l.CDN, "GET", URLInitQueryUserInfo, nil, h)
	if err != nil {
		Error("initQueryUserInfo DoForWardRequest error:", err)
	}
	Debug("initQueryUserInfo body:", body)
}

func (l *Login) leftTicketInit() {
	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/index/init"}
	body, err := DoForWardRequestHeader(l.CDN, "GET", URLInit, nil, h)
	if err != nil {
		Error("initQueryUserInfo DoForWardRequest error:", err)
	}
	Debug("initQueryUserInfo body:", body)
}
