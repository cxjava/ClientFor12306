package main

type Login struct {
	Username string
	Password string
	Captcha  string
	Cookie   string
}

type CheckRandCode struct {
	Basic
	Data             string        `json:"data"`
	Messages         []interface{} `json:"messages,omitempty"`
	ValidateMessages interface{}   `json:"validateMessages,omitempty"`
}

type LoginAysnSuggest struct {
	Basic
	Data struct {
		LoginCheck string `json:"loginCheck"`
	}
	Messages         []interface{} `json:"messages,omitempty"`
	ValidateMessages interface{}   `json:"validateMessages,omitempty"`
}
