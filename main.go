package main

import (
	"image"
	"log"
	"net/http"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	createWin()
}

func createWin() {
	img1 := GetImage("113.57.187.29")

	var mw *walk.MainWindow
	var acceptPB *walk.PushButton
	var iv *walk.ImageView

	bit, _ := walk.NewBitmapFromImage(img1)

	if _, err := (MainWindow{
		AssignTo: &mw,
		Title:    "Animal Details",
		MinSize:  Size{180, 210},
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
						AssignTo:    &iv,
						Image:       bit,
						MinSize:     Size{78, 26},
						MaxSize:     Size{78, 38},
						Name:        "captcha1",
						ToolTipText: "单击刷新验证码",
						OnMouseUp: func(x, y int, button walk.MouseButton) {
							img1 := GetImage("113.57.187.29")
							bit, _ = walk.NewBitmapFromImage(img1)
							iv.SetImage(bit)
						},
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

//获取新的验证码图片
func GetImage(cdn string) image.Image {
	req, err := http.NewRequest("GET", PassCodeNewURL, nil)
	if err != nil {
		Error("GetImage http.NewRequest error:", err)
		return nil
	}
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
	Debug("==" + GetCookieFromRespHeader(resp) + "==")
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
