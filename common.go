package main

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

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

//转发
func DoForWardRequest(cdn, method, requestUrl string, body io.Reader) (string, error) {
	return DoForWardRequestHeader(cdn, method, requestUrl, body, nil)
}

func DoForWardRequestHeader(cdn, method, requestUrl string, body io.Reader, customHeader map[string]string) (string, error) {
	req, err := http.NewRequest(method, requestUrl, body)
	if err != nil {
		Error("DoForWardRequest http.NewRequest error:", err)
		return "", err
	}
	AddReqestHeader(req, method, customHeader)

	con, err := NewForwardClientConn(cdn, req.URL.Scheme)
	if err != nil {
		Error("DoForWardRequest newForwardClientConn error:", err)
		return "", err
	}
	defer con.Close()
	resp, err := con.Do(req)
	if err != nil {
		Error("DoForWardRequest con.Do error:", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Error("DoForWardRequest StatusCode:", resp.StatusCode)
		for k, v := range resp.Header {
			Error("k=", k, "v=", v)
		}
		return "", err
	}
	content := ParseResponseBody(resp)
	Debug("DoForWardRequest content:", content)
	return content, nil
}

//添加头
func AddReqestHeader(request *http.Request, method string, customHeader map[string]string) {
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("Host", "kyfw.12306.cn")
	request.Header.Set("Accept-Encoding", "gzip,deflate")
	request.Header.Set("Accept-Language", "zh-CN")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	request.Header.Set("User-Agent", UserAgent)
	if login.Cookie != "" {
		request.Header.Set("Cookie", login.Cookie)
	}
	if method == "POST" {
		request.Header.Set("Content-Length", fmt.Sprintf("%d", request.ContentLength))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for k, v := range customHeader {
		request.Header.Set(k, v)
	}
	Info(request.Header)
}

//读取响应
func ParseResponseBody(resp *http.Response) string {
	var body []byte
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			Error("gzip.NewReader:", err)
			return ""
		}
		defer reader.Close()
		bodyByte, err := ioutil.ReadAll(reader)
		if err != nil {
			Error("gzip ioutil.ReadAll:", err)
			return ""
		}
		body = bodyByte
	default:
		bodyByte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Error("ioutil.ReadAll:", err)
			return ""
		}
		body = bodyByte
	}
	return string(body)
}

//	clientConn, _ := newForwardClientConn("www.google.com","https")
//	resp, err := clientConn.Do(req)
func NewForwardClientConn(forwardAddress, scheme string) (*httputil.ClientConn, error) {
	if "http" == scheme {
		conn, err := net.Dial("tcp", forwardAddress+":80")
		if err != nil {
			Error("newForwardClientConn net.Dial error:", err)
			return nil, err
		}
		return httputil.NewClientConn(conn, nil), nil
	}
	conn, err := tls.Dial("tcp", forwardAddress+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		Error("newForwardClientConn tls.Dial error:", err)
		return nil, err
	}
	return httputil.NewClientConn(conn, nil), nil
}

//获取车票余票信息
//getTicketNum("O008450822M010250252O008453240", "O0M0O0")
func GetTicketNum(yupiaoInfo, seat_types string) (ticketNum map[string]int) {
	ticketNum = make(map[string]int)
	//去除第一个类型，因为第一类型比较特殊，下面的str同样去掉
	types := strings.Split(seat_types[2:len(seat_types)-1], "0")
	//判断第一个类型
	if strings.HasPrefix(yupiaoInfo, "10") {
		num, _ := strconv.Atoi(yupiaoInfo[7:10])
		ticketNum["无座"] = num
	} else if strings.HasPrefix(yupiaoInfo, "O0") {
		num, _ := strconv.Atoi(yupiaoInfo[7:10])
		ticketNum["二等座"] = num
	} else if strings.HasPrefix(yupiaoInfo, "60") {
		num, _ := strconv.Atoi(yupiaoInfo[7:10])
		ticketNum["高级软卧"] = num
	} else {
		num, _ := strconv.Atoi(yupiaoInfo[7:10])
		ticketNum[yupiaoInfo[0:2]] = num
	}

	yupiaoInfo = yupiaoInfo[10:]

	for _, v := range types {
		key := v + "0"
		start := strings.Index(yupiaoInfo, key) + 7
		end := start + 3
		num, _ := strconv.Atoi(yupiaoInfo[start:end])
		switch key {
		case "10":
			ticketNum["硬座"] = num
		case "20":
			ticketNum["软座"] = num
		case "30":
			ticketNum["硬卧"] = num
		case "40":
			ticketNum["软卧"] = num
		case "O0":
			ticketNum["高铁无座"] = num
		case "M0":
			ticketNum["一等座"] = num
		case "90":
			ticketNum["商务座"] = num
		case "P0":
			ticketNum["特等座"] = num
		default:
			ticketNum[key] = num
		}
		yupiaoInfo = yupiaoInfo[end:]
	}
	return
}
