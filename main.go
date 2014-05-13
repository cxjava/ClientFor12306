package main

import (
	"encoding/json"
	"image"
	"log"
	"net/http"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	login  = &Login{}
	ticket = &TicketQueryInfo{}
)

type MyMainWindow struct {
	*walk.MainWindow
}

func main() {
	app := walk.App()

	// These specify the app data sub directory for the settings file.
	app.SetOrganizationName("The Walk Authors")
	app.SetProductName("Walk Settings Example")

	// Settings file name.
	settings := walk.NewIniFileSettings("settings.ini")

	// All settings marked as expiring will expire after this duration w/o use.
	// This applies to all widgets settings.
	settings.SetExpireDuration(time.Hour * 24 * 30 * 3)

	if err := settings.Load(); err != nil {
		log.Fatal(err)
	}

	app.SetSettings(settings)

	createTicketWin()
	// createLoginWin()

	if err := settings.Save(); err != nil {
		log.Fatal(err)
	}

}

var (
	mw           = &MyMainWindow{}
	loginButton  *walk.PushButton
	captchaImage *walk.ImageView
	captchaEdit  *walk.LineEdit
	loginDB      *walk.DataBinder
	loginEP      walk.ErrorPresenter

	ticketWin          = &MyMainWindow{}
	ticketDB           *walk.DataBinder
	ticketEP           walk.ErrorPresenter
	submitCaptchaImage *walk.ImageView
	submitCaptchaEdit  *walk.LineEdit
	seatTypeComboBox   *walk.ComboBox
)

func createTicketWin() {

	if _, err := (MainWindow{
		Name:     "ticketWindow",
		AssignTo: &ticketWin.MainWindow,
		Title:    "订票查询 -  by Charles",
		MinSize:  Size{300, 300},
		Layout:   VBox{},
		DataBinder: DataBinder{
			AssignTo:       &ticketDB,
			DataSource:     ticket,
			ErrorPresenter: ErrorPresenterRef{&ticketEP},
		},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 4},
				Name:   "ticketPanel",
				Children: []Widget{
					Label{
						Text: "出发日期:",
					},
					DateEdit{
						// MinDate: time.Now(),
						// MaxDate: time.Now().AddDate(0, 0, 20),
						Date: Bind("TrainDate"),
					},
					Label{
						Text: "席　别:",
					},
					ComboBox{
						AssignTo:      &seatTypeComboBox,
						Value:         Bind("SeatType", SelRequired{}),
						BindingMember: "Value",
						DisplayMember: "Name",
						Model:         KnownSeatTypeName(),
						OnCurrentIndexChanged: func() {
							Info(seatTypeComboBox.DisplayMember(), seatTypeComboBox.Name())
							Info(seatTypeComboBox.CurrentIndex(), seatTypeComboBox.Text(), seatTypeComboBox.BindingMember())
						},
					},

					Label{
						Text: "出发地:",
					},
					LineEdit{
						ColumnSpan: 3,
						MaxLength:  32,
						Text:       Bind("FromStationsStr"),
					},

					Label{
						Text: "目的地:",
					},
					LineEdit{
						ColumnSpan: 3,
						MaxLength:  32,
						Text:       Bind("ToStationsStr"),
					},

					Label{
						Text: "车　次:",
					},
					LineEdit{
						ColumnSpan: 3,
						MaxLength:  32,
						Text:       Bind("TriansStr"),
					},

					Label{
						Text: "验证码:",
					},
					LineEdit{
						AssignTo:  &submitCaptchaEdit,
						MaxLength: 4,
						OnKeyUp: func(key walk.Key) {
							if key == walk.KeyReturn && len(submitCaptchaEdit.Text()) == 4 {
								// mw.Submit()
							}
							// if len(captchaEdit.Text()) == 4 {
							// 	Info("no enter")
							// 	mw.Submit()
							// }
						},
					},
					VSpacer{
						ColumnSpan: 2,
						Size:       8,
					},
					VSpacer{
						ColumnSpan: 1,
						Size:       8,
					},
					ImageView{
						ColumnSpan: 1,
						AssignTo:   &submitCaptchaImage,
						// Image:       Im,
						MinSize:     Size{78, 38},
						MaxSize:     Size{78, 38},
						ToolTipText: "单击刷新验证码",
						OnMouseUp: func(x, y int, button walk.MouseButton) {
							i := GetImage(Conf.CDN[0])
							Im, _ := walk.NewBitmapFromImage(i)
							submitCaptchaImage.SetImage(Im)
							// seatTypeComboBox.SetCurrentIndex(6)
							// seatTypeComboBox.SetDisplayMember("3")
							// seatTypeComboBox.SetDisplayMember("硬卧")
							seatTypeComboBox.SetText("硬卧")
						},
					},
					LineErrorPresenter{
						AssignTo:   &ticketEP,
						ColumnSpan: 4,
					},
					VSpacer{
						ColumnSpan: 2,
						Size:       8,
					},
					VSpacer{
						ColumnSpan: 1,
						Size:       8,
					},
					PushButton{
						// AssignTo:  &loginButton,
						Text: "查询",
						// OnClicked: mw.Submit,
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}

}

func createLoginWin() {
	go func() {
		i := GetImage(Conf.CDN[0])
		Im, _ := walk.NewBitmapFromImage(i)
		captchaImage.SetImage(Im)
	}()

	if _, err := (MainWindow{
		Name:     "loginWindow",
		AssignTo: &mw.MainWindow,
		Title:    "登陆",
		MinSize:  Size{250, 250},
		Layout:   VBox{},
		DataBinder: DataBinder{
			AssignTo:       &loginDB,
			DataSource:     login,
			ErrorPresenter: ErrorPresenterRef{&loginEP},
		},

		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Name:   "loginPanel",
				Children: []Widget{
					Label{
						Text: "用户名:",
					},
					LineEdit{
						Name:      "Username",
						MaxLength: 30,
						Text:      Bind("Username"),
					},

					Label{
						Text: "密　码:",
					},
					LineEdit{
						Name:         "Password",
						MaxLength:    32,
						Text:         Bind("Password"),
						PasswordMode: true,
					},

					Label{
						Text: "验证码:",
					},
					LineEdit{
						AssignTo:  &captchaEdit,
						MaxLength: 4,
						Text:      Bind("Captcha"),
						OnKeyUp: func(key walk.Key) {
							if key == walk.KeyReturn && len(captchaEdit.Text()) == 4 {
								mw.Submit()
							}
							// if len(captchaEdit.Text()) == 4 {
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
						AssignTo: &captchaImage,
						// Image:       Im,
						MinSize:     Size{150, 60},
						MaxSize:     Size{150, 60},
						ToolTipText: "单击刷新验证码",
						OnMouseUp: func(x, y int, button walk.MouseButton) {
							i := GetImage(Conf.CDN[0])
							Im, _ := walk.NewBitmapFromImage(i)
							captchaImage.SetImage(Im)
						},
					},
					VSpacer{
						ColumnSpan: 1,
						Size:       8,
					},
					PushButton{
						AssignTo:  &loginButton,
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
	if err := loginDB.Submit(); err != nil {
		Error("login faild! :", err)
		return
	}
	Info(login)
	if r, m := login.CheckRandCodeAnsyn(Conf.CDN[0]); !r {
		msg := "验证码不正确！"
		if len(m) > 0 {
			msg = m[0]
		}
		captchaEdit.SetText("")
		captchaEdit.SetFocus()
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
		Im, _ := walk.NewBitmapFromImage(img)
		captchaImage.SetImage(Im)
		captchaEdit.SetText("")
		captchaEdit.SetFocus()
		return
	}
	mw.Dispose()
	createTicketWin()

}

//获取联系人
func getPassengerDTO() {
	passenger := &PassengerDTO{}
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
