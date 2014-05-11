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

type MyMainWindow struct {
	*walk.MainWindow
}

func main() {
	createWin()
}

var (
	mw       = new(MyMainWindow)
	submitPB *walk.PushButton
	iv       *walk.ImageView
	cp       *walk.LineEdit
	db       *walk.DataBinder
	ep       walk.ErrorPresenter
	bit      *walk.Bitmap
)

func createWin() {
	img1 := GetImage(Conf.CDN[0])

	bit, _ = walk.NewBitmapFromImage(img1)
	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
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
						AssignTo:  &cp,
						MaxLength: 4,
						Text:      Bind("Captcha"),
						OnKeyUp: func(key walk.Key) {
							if key == walk.KeyReturn && len(cp.Text()) == 4 {
								mw.Submit()
							}
							// if len(cp.Text()) == 4 {
							// 	Info("no enter")
							// 	mw.Submit()
							// }
						},
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
							img1 := GetImage(Conf.CDN[0])
							bit, _ = walk.NewBitmapFromImage(img1)
							iv.SetImage(bit)
						},
					},
					VSpacer{
						ColumnSpan: 1,
						Size:       8,
					},
					PushButton{
						AssignTo:  &submitPB,
						Text:      "登陆",
						OnClicked: mw.Submit,
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
	//set cookie
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
		Info("k=", k, "v=", v)
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
	crc := new(CheckRandCode)

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

	las := new(LoginAysnSuggest)
	if err := json.Unmarshal([]byte(content), &las); err != nil {
		Error("Login json.Unmarshal error:", err)
		return false, []string{err.Error()}
	}
	Info(las)
	return las.Data.LoginCheck == "Y", las.Messages
}
func (mw *MyMainWindow) Submit() {
	if err := db.Submit(); err != nil {
		Error("login faild! :", err)
		return
	}
	Info(login)
	if r, m := login.CheckRandCodeAnsyn(Conf.CDN[0]); !r {
		msg := "验证码不正确！"
		if len(m) > 0 {
			msg = m[0]
		}
		cp.SetText("")
		cp.SetFocus()
		walk.MsgBox(mw, "提示", msg, walk.MsgBoxIconInformation)
		return
	}
	if r, m := login.Login(Conf.CDN[0]); !r {
		msg := "系统错误！"
		if len(m) > 0 {
			msg = m[0]
		}
		walk.MsgBox(mw, "提示", msg, walk.MsgBoxIconInformation)
		img1 := GetImage(Conf.CDN[0])
		bit, _ = walk.NewBitmapFromImage(img1)
		iv.SetImage(bit)
		cp.SetText("")
		cp.SetFocus()
		return
	}
}
