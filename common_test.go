package main

import (
	"net/http"
	"strings"
	"testing"
)

func TestAddRequestHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "www.baidu.com", nil)
	if err != nil {
		t.Fatal("error:", err)
		return
	}
	AddReqestHeader(req, "GET")
	if req.Header.Get("Host") != "kyfw.12306.cn" {
		t.Fatal("AddReqestHeader failed!")
	}

	AddReqestHeader(req, "POST")
	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded; charset=UTF-8" {
		t.Fatal("AddReqestHeader for Content-Type failed!")
	}
}

func TestParseResponseBody(t *testing.T) {
	req, err := http.NewRequest("GET", "http://www.zhihu.com/read", nil)
	if err != nil {
		t.Fatal("error:", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("error:", err)
		return
	}
	content := ParseResponseBody(resp)
	if len(content) < 1 {
		t.Fatal("ParseResponseBody failed!")
	}
	if !strings.Contains(content, "知乎阅读") {
		t.Fatal("failded to get zhihu content!")
	}
}
func TestNewForwardClientConn(t *testing.T) {
	con, err := NewForwardClientConn("113.57.187.29", "https")
	if err != nil {
		t.Fatal("error:", err)
		return
	}
	if con == nil {
		t.Fatal("con is nil!")
	}

	con2, err := NewForwardClientConn("162.105.28.232", "http")
	if err != nil {
		t.Fatal("error:", err)
		return
	}
	if con2 == nil {
		t.Fatal("con is nil!")
	}
}

func TestGetTicketNum(t *testing.T) {
	ticketNum := GetTicketNum("O008450822M010250252O008453240", "O0M0O0")
	if len(ticketNum) < 1 {
		t.Fatal("GetTicketNum failed!")
	}
	if ticketNum["二等座"] != 822 {
		t.Fatal("GetTicketNum wrong number!")
	}
	if ticketNum["一等座"] != 252 {
		t.Fatal("GetTicketNum wrong number!")
	}
	if ticketNum["高铁无座"] != 240 {
		t.Fatal("GetTicketNum wrong number!")
	}
}
