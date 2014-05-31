package main

import (
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

var (
	passenger = &PassengerDTO{}
)

//获取联系人
func getPassengerDTO() {
	for _, cdn := range Conf.CDN {
		Info("开始获取联系人！")
		body, err := DoForWardRequest(cdn, "POST", URLGetPassengerDTOs, nil)
		if err != nil {
			Error("getPassengerDTO DoForWardRequest error:", err)
			continue
		}
		Debug("getPassengerDTO body:", body)

		if !strings.Contains(body, "passenger_name") {
			Error("获取联系人出错!!!!!!返回:", body)
			continue
		}

		if err := json.Unmarshal([]byte(body), &passenger); err != nil {
			Error("getPassengerDTO", cdn, err)
			continue
		} else {
			Info(cdn, "获取成功！")
			break
		}
	}
	Debug(passenger)

}

func dyLoginJs(cdn string) error {
	body, err := DoForWardRequest(cdn, "GET", URLLoginJs, nil)
	if err != nil {
		Error("dyLoginJs DoForWardRequest error:", err)
		return err
	}
	Debug("dyLoginJs body:", body)
	return nil
}

func dyQueryJs(cdn string) error {
	body, err := DoForWardRequest(cdn, "GET", URLQueryJs, nil)
	if err != nil {
		Error("dyQueryJs DoForWardRequest error:", err)
		return err
	}
	Debug("dyQueryJs body:", body)
	return nil
}

func initQueryUserInfo(cdn string) error {
	body, err := DoForWardRequest(cdn, "GET", URLInitQueryUserInfo, nil)
	if err != nil {
		Error("initQueryUserInfo DoForWardRequest error:", err)
		return err
	}
	Debug("initQueryUserInfo body:", body)
	return nil
}

func checkUser(cdn string) error {
	body, err := DoForWardRequest(cdn, "POST", URLCheckUser, nil)
	if err != nil {
		Error("checkUser DoForWardRequest error:", err)
		return err
	}
	Debug("checkUser body:", body)
	return nil
}

func leftTicketInit(cdn string) error {
	body, err := DoForWardRequest(cdn, "GET", URLInit, nil)
	if err != nil {
		Error("leftTicketInit DoForWardRequest error:", err)
		return err
	}
	Debug("leftTicketInit body:", body)
	return nil
}

//获取新的验证码图片
func GetLoginImage(cdn string) image.Image {
	req, err := http.NewRequest("GET", URLLoginPassCode, nil)
	if err != nil {
		Error("GetImage http.NewRequest error:", err)
		return nil
	}
	AddReqestHeader(req, "GET", map[string]string{"Accept": "image/webp,*/*;q=0.8"})
	req.Header.Del("X-Requested-With")
	con, err := NewForwardClientConn(cdn, req.URL.Scheme)
	if err != nil {
		Error("GetImage newForwardClientConn error:", err)
		return nil
	}
	defer con.Close()
	resp, err := con.Do(req)
	if err != nil {
		Error("GetImage con.Do error:", err)
		return nil
	}
	defer resp.Body.Close()

	img, s, err := image.Decode(resp.Body)
	Debug("image type:", s)
	if err != nil {
		Error("GetImage image.Decode:", err)
		return nil
	}
	return img
}

//从响应消息头里面获取cookie
func GetCookieFromRespHeader(resp *http.Response) (cookie string) {
	cookies := []string{}
	for k, v := range resp.Header {
		Debug("k=", k, "v=", v)
		if k == "Set-Cookie" {
			for _, b := range v {
				v := strings.Split(b, ";")[0]
				cookies = append(cookies, v)
				cookies = append(cookies, "; ")
			}
		}
	}
	d := strings.Join(cookies, "")
	if len(d) < 2 {
		return ""
	}
	cookie = d[:len(d)-2]
	return
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

	resp, err := client.Do(req)
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
	// request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("Accept-Language", "zh-CN")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("User-Agent", UserAgent)
	request.Header.Set("Accept-Encoding", "gzip,deflate")
	request.Header.Set("Host", "kyfw.12306.cn")
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
		return httputil.NewProxyClientConn(conn, nil), nil
	}
	conn, err := tls.Dial("tcp", forwardAddress+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		Error("newForwardClientConn tls.Dial error:", err)
		return nil, err
	}
	return httputil.NewProxyClientConn(conn, nil), nil
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
