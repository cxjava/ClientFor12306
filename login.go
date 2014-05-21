package main

import (
	"encoding/json"
	"net/url"
	"strings"
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

func CheckRandCodeAnsyn(randCode, cdn string) (r bool, msg []string) {
	b := url.Values{}
	b.Add("randCode", randCode)
	b.Add("rand", Rand)
	params, err := url.QueryUnescape(b.Encode())
	if err != nil {
		Error("CheckRandCodeAnsyn url.QueryUnescape error:", err)
		return false, []string{err.Error()}
	}
	Info(params)
	content, err := DoForWardRequest(cdn, "POST", URLCheckRandCodeAnsyn, strings.NewReader(params))
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

func (l *Login) Login(cdn string) (r bool, msg []string) {
	b := url.Values{}
	b.Add("loginUserDTO.user_name", l.Username)
	b.Add("userDTO.password", l.Password)
	b.Add("randCode", l.Captcha)
	params, err := url.QueryUnescape(b.Encode())
	if err != nil {
		Error("Login url.QueryUnescape error:", err)
		return false, []string{err.Error()}
	}
	content, err := DoForWardRequest(cdn, "POST", URLLoginAysnSuggest, strings.NewReader(params))
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
	body, err := DoForWardRequestHeader(Conf.CDN[0], "POST", URLUserLogin, strings.NewReader(params), h)
	if err != nil {
		Error("userLogin DoForWardRequest error:", err)
	}
	Debug("userLogin body:", body)
}
