package main

import (
	"encoding/json"
	"fmt"
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

var (
	login = new(Login)
)

type MyMainWindow struct {
	*walk.MainWindow
}

func main() {
	createLoginWin()
}

var (
	mw       = new(MyMainWindow)
	submitPB *walk.PushButton
	iv       *walk.ImageView
	cp       *walk.LineEdit
	db       *walk.DataBinder
	ep       walk.ErrorPresenter
	Im       *walk.Bitmap
)

func main2() {
	var mw *walk.MainWindow
	var outTE *walk.TextEdit

	if _, err := (MainWindow{
		AssignTo: &mw,
		Title:    "Walk Data Binding Example",
		MinSize:  Size{300, 200},
		Layout:   VBox{},
		Children: []Widget{
			PushButton{
				Text: "Edit Animal",
				OnClicked: func() {
					getPassengerDTO()
				},
			},
			Label{
				Text: "animal:",
			},
			TextEdit{
				AssignTo: &outTE,
				ReadOnly: true,
				Text:     fmt.Sprintf("%+v", login),
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

func createLoginWin() {
	i := GetImage(Conf.CDN[0])
	Im, _ = walk.NewBitmapFromImage(i)

	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "登陆",
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
						Image:       Im,
						MinSize:     Size{78, 26},
						MaxSize:     Size{78, 38},
						ToolTipText: "单击刷新验证码",
						OnMouseUp: func(x, y int, button walk.MouseButton) {
							i := GetImage(Conf.CDN[0])
							Im, _ = walk.NewBitmapFromImage(i)
							iv.SetImage(Im)
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

//登录逻辑
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
		img := GetImage(Conf.CDN[0])
		Im, _ = walk.NewBitmapFromImage(img)
		iv.SetImage(Im)
		cp.SetText("")
		cp.SetFocus()
		return
	}
	mw.Dispose()
	main2()

}

//获取联系人
func getPassengerDTO() {
	passenger := new(PassengerDTO)
	for _, cdn := range Conf.CDN {
		Info("开始获取联系人！")
		body, err := DoForWardRequest(cdn, "POST", GetPassengerDTOURL, nil)
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
	Info(passenger)
}
