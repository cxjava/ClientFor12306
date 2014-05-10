package main

import (
	"encoding/json"
	"image"
	"log"
	"net/http"
	"net/url"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	login = new(Login)
)

func main() {
	createWin()
}

func createWin() {
	img1 := GetImage("113.57.187.29")

	var mw *walk.MainWindow
	var submitPB *walk.PushButton
	var iv *walk.ImageView
	var db *walk.DataBinder
	var ep walk.ErrorPresenter

	bit, _ := walk.NewBitmapFromImage(img1)

	if _, err := (MainWindow{
		AssignTo: &mw,
		Title:    "Animal Details",
		MinSize:  Size{180, 210},
		Layout:   VBox{},
		DataBinder: DataBinder{
			AssignTo:       &db,
			DataSource:     login,
			ErrorPresenter: ErrorPresenterRef{&ep},
		},

		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "用户名:",
					},
					LineEdit{
						MaxLength: 30,
						Text:      Bind("Username"),
					},

					Label{
						Text: "密　码:",
					},
					LineEdit{
						MaxLength:    32,
						Text:         Bind("Password"),
						PasswordMode: true,
					},

					Label{
						Text: "验证码:",
					},
					LineEdit{
						MaxLength: 4,
						Text:      Bind("Captcha"),
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
						AssignTo: &submitPB,
						Text:     "登陆",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								Error("login faild! :", err)
								return
							}
							Info(login)
							Info(login.CheckRandCodeAnsyn("113.57.187.29"))
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

	login.Cookie = GetCookieFromRespHeader(resp)
	Debug("==" + login.Cookie + "==")

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

func (l *Login) CheckRandCodeAnsyn(cdn string) bool {
	b := url.Values{}
	b.Add("randCode", l.Captcha)
	b.Add("rand", Rand)
	params, err := url.QueryUnescape(b.Encode())
	if err != nil {
		Error("CheckRandCodeAnsyn url.QueryUnescape error:", err)
		return false
	}

	req, err := http.NewRequest("POST", CheckRandCodeURL, strings.NewReader(params))
	if err != nil {
		Error("CheckRandCodeAnsyn http.NewRequest error:", err)
		return false
	}
	AddReqestHeader(req, "POST")

	con, err := NewForwardClientConn(cdn, req.URL.Scheme)
	if err != nil {
		Error("CheckRandCodeAnsyn newForwardClientConn error:", err)
		return false
	}
	defer con.Close()
	resp, err := con.Do(req)
	if err != nil {
		Error("CheckRandCodeAnsyn con.Do error:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Error("CheckRandCodeAnsyn StatusCode:", resp.StatusCode, resp.Header, resp.Cookies())
		return false
	}
	content := ParseResponseBody(resp)
	Debug("CheckRandCodeAnsyn content:", content)

	crc := new(CheckRandCode)

	if err := json.Unmarshal([]byte(content), &crc); err != nil {
		Error("CheckRandCodeAnsyn json.Unmarshal error:", err)
		return false
	}
	Info(crc)
	return crc.Data == "Y"
}
