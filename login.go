package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type UserLoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Code     string `form:"code" binding:"required"`
}

type Login struct {
	Username string
	Password string
	Captcha  string
	Cookie   string
	CDN      string
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

	// con, err := NewForwardClientConn(l.CDN, req.URL.Scheme)
	// if err != nil {
	// 	Error("DoForWardRequestHeader NewForwardClientConn error:", err)
	// 	return "", err
	// }
	// defer con.Close()
	// resp, err := con.Do(req)
	//
	resp, err := client.Do(req)
	if err != nil {
		Error("setNewCookie error:", err)
		return err
	}
	defer resp.Body.Close()

	l.Cookie = GetCookieFromRespHeader(resp)
	Info("Get New Cookie=" + l.Cookie + "=")
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
	Info("CheckRandCodeAnsyn params:", b.Encode())
	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/login/init"}
	content, err := DoForWardRequestHeader(login.CDN, "POST", URLCheckRandCodeAnsyn, strings.NewReader(b.Encode()), h)
	if err != nil {
		Error("CheckRandCodeAnsyn DoForWardRequestHeader error:", err)
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
	h := map[string]string{
		"X-Requested-With": "XMLHttpRequest",
		"Referer":          "https://kyfw.12306.cn/otn/leftTicket/init",
	}
	body, err := DoForWardRequestHeader(l.CDN, "POST", URLCheckUser, strings.NewReader("_json_att="), h)
	if err != nil {
		Error("checkUser DoForWardRequest error:", err)
		return false, err
	}
	Debug("checkUser body:", body)
	return strings.Contains(body, `"data":{"flag":true}`), nil
}

func (l *Login) loginAysnSuggest() (r bool, msg []string) {
	b := url.Values{}
	b.Add("loginUserDTO.user_name", l.Username)
	b.Add("userDTO.password", l.Password)
	b.Add("randCode", l.Captcha)
	params, err := url.QueryUnescape(b.Encode())
	if err != nil {
		Error("loginAysnSuggest url.QueryUnescape error:", err)
		return false, []string{err.Error()}
	}
	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/login/init"}
	content, err := DoForWardRequestHeader(l.CDN, "POST", URLLoginAysnSuggest, strings.NewReader(params), h)
	if err != nil {
		Error("loginAysnSuggest DoForWardRequest error:", err)
		return false, []string{err.Error()}
	}

	las := &LoginAysnSuggest{}
	if err := json.Unmarshal([]byte(content), &las); err != nil {
		Error("loginAysnSuggest json.Unmarshal error:", err)
		return false, []string{err.Error()}
	}
	Info(las)
	return las.Data.LoginCheck == "Y", las.Messages
}

func (l *Login) userLogin() {
	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/login/init"}
	body, err := DoForWardRequestHeader(l.CDN, "POST", URLUserLogin, strings.NewReader("_json_att="), h)
	if err != nil {
		Error("userLogin DoForWardRequestHeader error:", err)
	}
	Debug("userLogin body:", body)
}

//获取联系人
func (l *Login) getPassengerDTO() (p *PassengerDTO) {
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
	p = &PassengerDTO{}
	if err := json.Unmarshal([]byte(body), &p); err != nil {
		Error("getPassengerDTO", l.CDN, err)
		return
	}
	return
}

func (l *Login) initQueryUserInfo() {
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
		Error("leftTicketInit DoForWardRequestHeader error:", err)
	}
	Debug("leftTicketInit body:", body)
}
