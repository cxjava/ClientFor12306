package main

import (
	"encoding/json"
	"image"
	"log"
	"math/rand"
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
	ticket = &TicketQueryInfo{
		P1: &PassengerOrder{},
		P2: &PassengerOrder{},
		P3: &PassengerOrder{},
		P4: &PassengerOrder{},
		P5: &PassengerOrder{},
	}
	passenger     = &PassengerDTO{}
	ticketWin     = &MyMainWindow{}
	loginWin      = &MyMainWindow{}
	myPassengers  = &walk.ComboBox{}
	mapPassengers = make(map[string]Passenger)

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
	date               *walk.DateEdit
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

	createLoginWin()
	createTicketWin()

	if err := settings.Save(); err != nil {
		log.Fatal(err)
	}

}

func createTicketWin() {
	go getPassengerDTO()
	go func() {
		time.Sleep(time.Second * 2)
		date.SetRange(time.Now(), time.Now().AddDate(0, 0, 19))
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
							if key == walk.KeyReturn && len(submitCaptchaEdit.Text()) == 4 {
								// loginWin.Submit()
							}
							// submitCaptchaEdit1.SetWidth(120)
							// if len(captchaEdit.Text()) == 4 {
							// 	Info("no enter")
							// 	loginWin.Submit()
							// }
						},
					},
					ImageView{
						AssignTo: &submitCaptchaImage,
						// Image:       Im,
						MinSize:     Size{78, 26},
						MaxSize:     Size{78, 26},
						ToolTipText: "单击刷新验证码",
						OnMouseUp: func(x, y int, button walk.MouseButton) {
							i := GetImage(Conf.CDN[0])
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
							ticket.queryLeftTicket(Conf.CDN[0])
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

//查询余票
func (t *TicketQueryInfo) queryLeftTicket(cdn string) *QueryLeftNewDTO {
	fr := t.FromStations
	to := t.ToStations
	Info(fr)
	Info(to)
	Info(len(fr))
	Info(rand.Intn(len(fr)))
	return nil
	leftTicketUrl := QueryLeftTicketURL

	leftTicketUrl += "leftTicketDTO.train_date=" + t.TrainDate.Format("2006-01-02") + "&"
	leftTicketUrl += "leftTicketDTO.from_station=" + StationMap[fr[rand.Intn(len(fr))]] + "&"
	leftTicketUrl += "leftTicketDTO.to_station=" + StationMap[to[rand.Intn(len(to))]] + "&"
	leftTicketUrl += "purpose_codes=ADULT"

	Debug("queryLeftTicket url:", leftTicketUrl)

	Info("开始获取联系人！")
	body, err := DoForWardRequest(cdn, "POST", leftTicketUrl, nil)
	if err != nil {
		Error("queryLeftTicket DoForWardRequest error:", err)
		return nil
	}
	Debug("queryLeftTicket body:", body)

	if !strings.Contains(body, "queryLeftNewDTO") {
		Error("查询余票出错，返回:", body, "查询链接:", leftTicketUrl)
		//删除废弃的CDN
		// if len(availableCDN) > 5 {
		// delete(availableCDN, cdn)
		// }
		return nil
	}
	leftTicket := &QueryLeftNewDTO{}

	if err := json.Unmarshal([]byte(body), &leftTicket); err != nil {
		Error("queryLeftTicket", cdn, err)
		return nil
	} else {
		Info(cdn, "获取成功！")
	}

	return leftTicket
}

func parseTicket() {
	Info(ticket.TrainDate.Format("2006-01-02"))
	Info(ticket)
	ticket.FromStations = parseStrings(ticket.FromStationsStr)
	ticket.ToStations = parseStrings(ticket.ToStationsStr)
	ticket.Trians = parseStrings(ticket.TriansStr)

	Info(ticket)
	o, n := parseStranger(*ticket)
	ticket.OldPassengerStr = o
	ticket.PassengerTicketStr = n[:len(n)-1]
	Info(ticket)
}
func parseStranger(ticket TicketQueryInfo) (oStr, nStr string) {
	if strings.Trim(ticket.P1.Name, " ") != "" {
		pa := ticket.P1
		nStr += pa.SeatType + ",0," + pa.TicketType + "," + pa.Name + "," + pa.PassengerIdTypeCode + "," + pa.PassengerIdNo + ",,N_"
		oStr += pa.Name + "," + pa.PassengerIdTypeCode + "," + pa.PassengerIdNo + "," + pa.TicketType + "_"
	}
	if strings.Trim(ticket.P2.Name, " ") != "" {
		pa := ticket.P2
		nStr += pa.SeatType + ",0," + pa.TicketType + "," + pa.Name + "," + pa.PassengerIdTypeCode + "," + pa.PassengerIdNo + ",,N_"
		oStr += pa.Name + "," + pa.PassengerIdTypeCode + "," + pa.PassengerIdNo + "," + pa.TicketType + "_"
	}
	if strings.Trim(ticket.P3.Name, " ") != "" {
		pa := ticket.P3
		nStr += pa.SeatType + ",0," + pa.TicketType + "," + pa.Name + "," + pa.PassengerIdTypeCode + "," + pa.PassengerIdNo + ",,N_"
		oStr += pa.Name + "," + pa.PassengerIdTypeCode + "," + pa.PassengerIdNo + "," + pa.TicketType + "_"
	}
	if strings.Trim(ticket.P4.Name, " ") != "" {
		pa := ticket.P4
		nStr += pa.SeatType + ",0," + pa.TicketType + "," + pa.Name + "," + pa.PassengerIdTypeCode + "," + pa.PassengerIdNo + ",,N_"
		oStr += pa.Name + "," + pa.PassengerIdTypeCode + "," + pa.PassengerIdNo + "," + pa.TicketType + "_"
	}
	if strings.Trim(ticket.P5.Name, " ") != "" {
		pa := ticket.P5
		nStr += pa.SeatType + ",0," + pa.TicketType + "," + pa.Name + "," + pa.PassengerIdTypeCode + "," + pa.PassengerIdNo + ",,N_"
		oStr += pa.Name + "," + pa.PassengerIdTypeCode + "," + pa.PassengerIdNo + "," + pa.TicketType + "_"
	}
	return
}
func parseStrings(str string) (s []string) {
	if strings.ContainsRune(str, rune('，')) {
		for _, v := range strings.Split(str, "，") {
			if v != "" {
				s = append(s, v)
			}
		}
	}
	if strings.ContainsRune(str, rune(',')) {
		for _, v := range strings.Split(str, ",") {
			if v != "" {
				s = append(s, v)
			}
		}
	}
	return
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
func createLoginWin() {
	go func() {
		i := GetImage(Conf.CDN[0])
		Im, _ := walk.NewBitmapFromImage(i)
		captchaImage.SetImage(Im)

		username.SetText("xuhong157499")
		password.SetText("xuhong1990")
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
						OnClicked: loginWin.Submit,
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

//登录逻辑
func (loginWin *MyMainWindow) Submit() {
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
		walk.MsgBox(loginWin, "提示", msg, walk.MsgBoxIconInformation)
		return
	}
	if r, m := login.Login(Conf.CDN[0]); !r {
		msg := "系统错误！"
		if len(m) > 0 {
			msg = m[0]
		}
		walk.MsgBox(loginWin, "提示", msg, walk.MsgBoxIconInformation)
		img := GetImage(Conf.CDN[0])
		Im, _ := walk.NewBitmapFromImage(img)
		captchaImage.SetImage(Im)
		captchaEdit.SetText("")
		captchaEdit.SetFocus()
		return
	}
	Info("登录成功！")
	loginWin.Dispose()
}

//获取联系人
func getPassengerDTO() {
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
	Debug(passenger)

	go func() {
		model := []string{}
		for _, v1 := range passenger.Data.NormalPassengers {
			model = append(model, v1.PassengerName)
			mapPassengers[v1.PassengerName] = v1
		}
		myPassengers.SetModel(model)
	}()
}
