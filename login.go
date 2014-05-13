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

func (l *Login) CheckRandCodeAnsyn(cdn string) (r bool, msg []string) {
	b := url.Values{}
	b.Add("randCode", l.Captcha)
	b.Add("rand", Rand)
	params, err := url.QueryUnescape(b.Encode())
	if err != nil {
		Error("CheckRandCodeAnsyn url.QueryUnescape error:", err)
		return false, []string{err.Error()}
	}
	content, err := DoForWardRequest(cdn, "POST", CheckRandCodeURL, strings.NewReader(params))
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
	content, err := DoForWardRequest(cdn, "POST", LoginAysnSuggestURL, strings.NewReader(params))
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
