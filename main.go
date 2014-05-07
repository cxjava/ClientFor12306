package main

import (
	"crypto/tls"
	"fmt"
	"image"
	"log"
	"net"
	"net/http"
	"net/http/httputil"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {

	req, err := http.NewRequest("GET", "https://kyfw.12306.cn/otn/passcodeNew/getPassCodeNew?module=login&rand=sjrand", nil)
	if err != nil {
		fmt.Println("doRequest http.NewRequest error:", err)
		return
	}
	clientConn, err := newForwardClientConn("113.57.187.29", req.URL.Scheme)
	if err != nil {
		fmt.Println("doRequest newForwardClientConn error:", err)
		return
	}
	defer clientConn.Close()
	resp, err := clientConn.Do(req)

	// resp, err := http.Get("http://kyfw.12306.cn/otn/passcodeNew/getPassCodeNew?module=login&rand=sjrand")
	if err != nil {
		fmt.Println("Client.Do:", err)
		return
	}
	defer resp.Body.Close()
	img1, s, err := image.Decode(resp.Body)
	fmt.Println("Decode.Do:", s)
	if err != nil {
		fmt.Println("Decode.Do:", err)
		return
	}

	var mw *walk.MainWindow
	var acceptPB *walk.PushButton
	bit, _ := walk.NewBitmapFromImage(img1)
	// var imageView *walk.ImageView
	// imageView.SetImage(bit)

	if _, err := (MainWindow{
		AssignTo: &mw,
		Title:    "Animal Details",
		MinSize:  Size{180, 220},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "用户名:",
					},
					LineEdit{
						Name: "username",
					},

					Label{
						Text: "密　码:",
					},
					LineEdit{
						Name:         "password",
						PasswordMode: true,
					},

					Label{
						Text: "验证码:",
					},
					LineEdit{
						Name: "captcha",
					},
					VSpacer{
						ColumnSpan: 1,
						Size:       8,
					},
					ImageView{
						Image:   bit,
						MinSize: Size{78, 26},
						MaxSize: Size{78, 38},
						// ColumnSpan: 2,
						Name: "captcha1",
					},
					VSpacer{
						ColumnSpan: 1,
						Size:       8,
					},
					PushButton{
						// ColumnSpan: 2,
						AssignTo: &acceptPB,
						Text:     "登陆",
						OnClicked: func() {

						},
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

//	clientConn, _ := newForwardClientConn("www.google.com","https")
//	resp, err := clientConn.Do(req)
func newForwardClientConn(forwardAddress, scheme string) (*httputil.ClientConn, error) {
	if "http" == scheme {
		conn, err := net.Dial("tcp", forwardAddress+":80")
		if err != nil {
			fmt.Println("newForwardClientConn net.Dial error:", err)
			return nil, err
		}
		return httputil.NewClientConn(conn, nil), nil
	}
	conn, err := tls.Dial("tcp", forwardAddress+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println("newForwardClientConn tls.Dial error:", err)
		return nil, err
	}
	return httputil.NewClientConn(conn, nil), nil
}
