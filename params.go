package main

type header map[string]string

var HeaderMap = map[string]header{
	"https://kyfw.12306.cn/otn/passcodeNew/checkRandCodeAnsyn": {
		"Accept":           "*/*",
		"Origin":           "https://kyfw.12306.cn",
		"X-Requested-With": "XMLHttpRequest",
	},
	"https://kyfw.12306.cn/otn/login/loginAysnSuggest": {
		"Accept":           "*/*",
		"Origin":           "https://kyfw.12306.cn",
		"X-Requested-With": "XMLHttpRequest",
		"Referer":          "https://kyfw.12306.cn/otn/login/init",
	},
	"https://kyfw.12306.cn/otn/login/userLogin": {
		"Accept":           "*/*",
		"Origin":           "https://kyfw.12306.cn",
		"X-Requested-With": "XMLHttpRequest",
		"Referer":          "https://kyfw.12306.cn/otn/login/init",
	},
	"https://kyfw.12306.cn/otn/login/checkUser": {
		"Accept":           "*/*",
		"Origin":           "https://kyfw.12306.cn",
		"X-Requested-With": "XMLHttpRequest",
	},
	"https://kyfw.12306.cn/otn/passcodeNew/getPassCodeNew.do?module=login&rand=sjrand&0.2680066132452339": {
		"Accept":  "image/webp,*/*;q=0.8",
		"Referer": "https://kyfw.12306.cn/otn/login/init",
	},
	"https://kyfw.12306.cn/otn/leftTicket/query?": {
		"Accept":           "*/*",
		"Cache-Control":    "no-store,no-cache",
		"Pragma":           "no-cache",
		"X-Requested-With": "XMLHttpRequest",
	}, "https://kyfw.12306.cn/otn/dynamicJs/queryJs": {
		"Accept":        "*/*",
		"Cache-Control": "no-cache",
	},
	"https://kyfw.12306.cn/otn/leftTicket/log?": {
		"Accept":           "*/*",
		"Cache-Control":    "no-store,no-cache",
		"Pragma":           "no-cache",
		"X-Requested-With": "XMLHttpRequest",
	},
	"https://kyfw.12306.cn/otn/confirmPassenger/getPassengerDTOs": {
		"Accept":           "*/*",
		"Origin":           "https://kyfw.12306.cn",
		"X-Requested-With": "XMLHttpRequest",
	},
	"https://kyfw.12306.cn/otn/confirmPassenger/autoSubmitOrderRequest": {
		"Accept":           "*/*",
		"Origin":           "https://kyfw.12306.cn",
		"X-Requested-With": "XMLHttpRequest",
	},
	"https://kyfw.12306.cn/otn/confirmPassenger/confirmSingle": {
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"Origin":           "https://kyfw.12306.cn",
		"X-Requested-With": "XMLHttpRequest",
	},
	"https://kyfw.12306.cn/otn/confirmPassenger/getQueueCountAsync": {
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"Origin":           "https://kyfw.12306.cn",
		"X-Requested-With": "XMLHttpRequest",
	},
}
