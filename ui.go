package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	ticketWin    = &MyMainWindow{}
	loginWin     = &MyMainWindow{}
	myPassengers = &walk.ComboBox{}

	loginButton  *walk.PushButton
	captchaImage *walk.ImageView
	captchaEdit  *walk.LineEdit
	loginDB      *walk.DataBinder
	loginEP      walk.ErrorPresenter

	ticketDB           *walk.DataBinder
	ticketEP           walk.ErrorPresenter
	submitCaptchaImage *walk.ImageView
	username           *walk.LineEdit
	password           *walk.LineEdit
	submitCaptchaEdit  *walk.LineEdit
	submitCaptchaEdit1 *walk.LineEdit
	passengers         *walk.Composite
	date               = &walk.DateEdit{}
)

type MyMainWindow struct {
	*walk.MainWindow
}

func setSubmitImage() {
	i := GetImage(Conf.CDN[0], false)
	Im, _ := walk.NewBitmapFromImage(i)
	submitCaptchaImage.SetImage(Im)
	submitCaptchaEdit.SetText("")
	submitCaptchaEdit.SetFocus()
}
func createUI() {
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

	createLoginWin()

	if err := settings.Save(); err != nil {
		log.Fatal(err)
	}

}

func createTicketWin() {
	go func() {
		getPassengerDTO()
		model := []string{}
		for _, v1 := range passenger.Data.NormalPassengers {
			model = append(model, v1.PassengerName)
			mapPassengers[v1.PassengerName] = v1
		}
		myPassengers.SetModel(model)
	}()

	go func() {
		time.Sleep(time.Second * 1)
		date.SetRange(time.Now(), time.Now().AddDate(0, 0, 19))
		date.SetDate(time.Now().AddDate(0, 0, 19))
	}()
	if _, err := (MainWindow{
		Name:     "ticketWindow",
		AssignTo: &ticketWin.MainWindow,
		Title:    "订票查询 -  by Charles",
		MinSize:  Size{300, 200},
		Layout:   VBox{},
		DataBinder: DataBinder{
			AssignTo:       &ticketDB,
			DataSource:     ticket,
			ErrorPresenter: ErrorPresenterRef{&ticketEP},
		},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 5},
				Name:   "ticketPanel",
				Children: []Widget{
					Label{
						Text: "出发日期:",
					},
					DateEdit{
						AssignTo: &date,
						Date:     Bind("TrainDate", SelRequired{}),
					},
					Label{
						Text: "车　次:",
					},
					LineEdit{
						ToolTipText: "多个车次请以逗号分隔",
						MaxLength:   32,
						Text:        Bind("TriansStr", SelRequired{}),
					},
					VSpacer{
						Size: 8,
					},

					Label{
						Text: "出发地:",
					},
					LineEdit{
						ToolTipText: "多个出发地请以逗号分隔",
						MaxLength:   32,
						Text:        Bind("FromStationsStr", SelRequired{}),
					},

					Label{
						Text: "目的地:",
					},
					LineEdit{
						MaxLength:   32,
						ToolTipText: "多个目的地请以逗号分隔",
						Text:        Bind("ToStationsStr", SelRequired{}),
					},
					ComboBox{
						AssignTo:              &myPassengers,
						BindingMember:         "Value",
						DisplayMember:         "Name",
						ToolTipText:           "选择联系人",
						OnCurrentIndexChanged: choosePassengers,
					},
				},
			},
			Composite{
				AssignTo: &passengers,
				Layout:   Grid{Columns: 5},
				Children: []Widget{
					LineEdit{
						Text: Bind("P1.Name"),
					},
					ComboBox{
						CurrentIndex:  0,
						Value:         Bind("P1.TicketType"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         KnownTicketTypeName(),
					},
					ComboBox{
						CurrentIndex:  0,
						Value:         Bind("P1.PassengerIdTypeCode"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         KnownIDTypeName(),
					},
					LineEdit{
						MinSize: Size{140, 12},
						Text:    Bind("P1.PassengerIdNo"),
					},
					ComboBox{
						Value:         Bind("P1.SeatType"),
						BindingMember: "Value",
						DisplayMember: "Name",
						Model:         KnownSeatTypeName(),
					},

					LineEdit{
						Text: Bind("P2.Name"),
					},
					ComboBox{
						Value:         Bind("P2.TicketType"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         KnownTicketTypeName(),
					},
					ComboBox{
						Value:         Bind("P2.PassengerIdTypeCode"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         KnownIDTypeName(),
					},
					LineEdit{
						Text: Bind("P2.PassengerIdNo"),
					},
					ComboBox{
						Value:         Bind("P2.SeatType"),
						BindingMember: "Value",
						DisplayMember: "Name",
						Model:         KnownSeatTypeName(),
					},

					LineEdit{
						Text: Bind("P3.Name"),
					},
					ComboBox{
						Value:         Bind("P3.TicketType"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         KnownTicketTypeName(),
					},
					ComboBox{
						Value:         Bind("P3.PassengerIdTypeCode"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         KnownIDTypeName(),
					},
					LineEdit{
						Text: Bind("P3.PassengerIdNo"),
					},
					ComboBox{
						Value:         Bind("P3.SeatType"),
						BindingMember: "Value",
						DisplayMember: "Name",
						Model:         KnownSeatTypeName(),
					},

					LineEdit{
						Text: Bind("P4.Name"),
					},
					ComboBox{
						Value:         Bind("P4.TicketType"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         KnownTicketTypeName(),
					},
					ComboBox{
						Value:         Bind("P4.PassengerIdTypeCode"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         KnownIDTypeName(),
					},
					LineEdit{
						Text: Bind("P4.PassengerIdNo"),
					},
					ComboBox{
						Value:         Bind("P4.SeatType"),
						BindingMember: "Value",
						DisplayMember: "Name",
						Model:         KnownSeatTypeName(),
					},

					LineEdit{
						Text: Bind("P5.Name"),
					},
					ComboBox{
						Value:         Bind("P5.TicketType"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         KnownTicketTypeName(),
					},
					ComboBox{
						Value:         Bind("P5.PassengerIdTypeCode"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         KnownIDTypeName(),
					},
					LineEdit{
						Text: Bind("P5.PassengerIdNo"),
					},
					ComboBox{
						Value:         Bind("P5.SeatType"),
						BindingMember: "Value",
						DisplayMember: "Name",
						Model:         KnownSeatTypeName(),
					},
				},
			},
			Composite{
				// Layout: Grid{Columns: 5},
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "验证码:",
					},
					LineEdit{
						AssignTo:  &submitCaptchaEdit,
						MaxLength: 4,
						OnKeyUp: func(key walk.Key) {
							if len(submitCaptchaEdit.Text()) == 4 {
								/*if r, m := order.checkRandCodeAnsyn(submitCaptchaEdit.Text()); !r {
									msg := "验证码不正确！"
									Info(msg)
									if len(m) > 0 {
										msg = m[0]
									}
									submitCaptchaEdit.SetText("")
									submitCaptchaEdit.SetFocus()
									walk.MsgBox(ticketWin, "提示", msg, walk.MsgBoxIconInformation)
									return
								}*/
								Info("success!")
								order.RandCode = submitCaptchaEdit.Text()
								go order.checkOrderInfo()
							}
							if key == walk.KeyReturn && len(submitCaptchaEdit.Text()) == 4 {
							}
							Info("over")
						},
					},
					ImageView{
						AssignTo: &submitCaptchaImage,
						// Image:       Im,
						MinSize:     Size{78, 26},
						MaxSize:     Size{78, 26},
						ToolTipText: "单击刷新验证码",
						OnMouseUp: func(x, y int, button walk.MouseButton) {
							i := GetImage(Conf.CDN[0], false)
							Im, _ := walk.NewBitmapFromImage(i)
							submitCaptchaImage.SetImage(Im)
						},
					},
					PushButton{
						Text: "查询",
						OnClicked: func() {
							if err := ticketDB.Submit(); err != nil {
								Error("login faild! :", err)
								return
							}
							parseTicket()
						},
					},
					LineErrorPresenter{
						AssignTo: &ticketEP,
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
		i := GetImage(Conf.CDN[0], true)
		Im, _ := walk.NewBitmapFromImage(i)
		captchaImage.SetImage(Im)

		username.SetText("cxjava11")
		password.SetText("Kee2209D6a050e5E")

		fmt.Println(strconv.FormatInt(time.Now().UnixNano(), 10))
		fmt.Println(time.Now().Unix())
		Info(time.Now().UTC().Local().Format(time.RFC1123Z))
		Info(time.Now().Local().Format(`Mon Jan 02 2006 15:04:05 GMT-0700 (China Standard Time)`))
	}()

	if _, err := (MainWindow{
		Name:     "loginWindow",
		AssignTo: &loginWin.MainWindow,
		Title:    "登陆",
		MinSize:  Size{70, 70},
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
						AssignTo:  &username,
						Name:      "Username",
						MaxLength: 30,
						Text:      Bind("Username"),
					},

					Label{
						Text: "密　码:",
					},
					LineEdit{
						AssignTo:     &password,
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
								loginWin.Submit()
							}
							// if len(captchaEdit.Text()) == 4 {
							// 	Info("no enter")
							// 	loginWin.Submit()
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
						// MinSize:     Size{150, 60},
						// MaxSize:     Size{150, 60},
						MinSize:     Size{78, 26},
						MaxSize:     Size{78, 26},
						ToolTipText: "单击刷新验证码",
						OnMouseUp: func(x, y int, button walk.MouseButton) {
							i := GetImage(Conf.CDN[0], true)
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
						OnClicked: loginWin.Submit,
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

//选择联系人
func choosePassengers() {
	p := mapPassengers[myPassengers.Text()]
	c := passengers.Children()
	for i := 0; i < 5; i++ {
		name, _ := c.At(i*5 + 0).(*walk.LineEdit)
		// Info(name.Text())
		if strings.Trim(name.Text(), " ") == "" {
			name.SetText(p.PassengerName)

			ticketType, _ := c.At(i*5 + 1).(*walk.ComboBox)
			ticketType.SetCurrentIndex(0)

			noType, _ := c.At(i*5 + 2).(*walk.ComboBox)
			noType.SetCurrentIndex(0)

			IDNO, _ := c.At(i*5 + 3).(*walk.LineEdit)
			// Info(IDNO.Text())
			IDNO.SetText(p.PassengerIdNo)

			seatType, _ := c.At(i*5 + 4).(*walk.ComboBox)
			seatType.SetCurrentIndex(6)
			break
		}
	}
}

//登录逻辑
func (loginWin *MyMainWindow) Submit() {
	if err := loginDB.Submit(); err != nil {
		Error("login faild! :", err)
		return
	}
	Info(login)
	if r, m := CheckRandCodeAnsyn(captchaEdit.Text(), Conf.CDN[0]); !r {
		msg := "验证码不正确！"
		if len(m) > 0 {
			msg = m[0]
		}
		captchaEdit.SetText("")
		captchaEdit.SetFocus()
		walk.MsgBox(loginWin, "提示", msg, walk.MsgBoxIconInformation)
		return
	}
	go DoForWardRequest(Conf.CDN[0], "POST", CheckUserURL, nil)
	if r, m := login.Login(Conf.CDN[0]); !r {
		msg := "系统错误！"
		if len(m) > 0 {
			msg = m[0]
		}
		walk.MsgBox(loginWin, "提示", msg, walk.MsgBoxIconInformation)
		img := GetImage(Conf.CDN[0], true)
		Im, _ := walk.NewBitmapFromImage(img)
		captchaImage.SetImage(Im)
		captchaEdit.SetText("")
		captchaEdit.SetFocus()
		return
	}
	go DoForWardRequest(Conf.CDN[0], "POST", CheckUserURL, nil)
	go loginCheck(Conf.CDN[0])
	Info("登录成功！")
	loginWin.Dispose()
	createTicketWin()
}
