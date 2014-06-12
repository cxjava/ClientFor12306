package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

var (
	client = &http.Client{}
	ws     *websocket.Conn
)

func init() {
	t := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: func(network, addr string) (net.Conn, error) {
			deadline := time.Now().Add(10 * time.Second)
			c, err := net.DialTimeout(network, addr, 10*time.Second)
			// c, err := net.DialTimeout(network, Conf.CDN[0]+":443", 10*time.Second)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(deadline)
			return c, nil
		},
	}
	if Conf.Proxy {
		pr, err := url.Parse(Conf.ProxyUrl)
		if err != nil {
			Error(err)
			return
		}
		t.Proxy = http.ProxyURL(pr)
	}
	client = &http.Client{
		Transport: t,
	}

}

type UserLoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Code     string `form:"code" binding:"required"`
}
type TicketQuery struct {
	CDN                string
	FromStations       []string
	ToStations         []string
	Trians             []string
	PassengerTicketStr string
	OldPassengerStr    string
	NumOfSeatType      map[string]int
	Start              string `form:"start" binding:"required"`
	End                string `form:"end" binding:"required"`
	Train              string `form:"train" binding:"required"`
	TrainDate          string `form:"date" binding:"required"`
	P1                 struct {
		PassengerName1 string `form:"passengerName1" binding:"required"`
		TicketType1    string `form:"ticketType1" binding:"required"`
		SeatType1      string `form:"seatType1" binding:"required"`
		IDType1        string `form:"IDType1" binding:"required"`
		IDNumber1      string `form:"IDNumber1" binding:"required"`
	}
	P2 struct {
		PassengerName2 string `form:"passengerName2" `
		TicketType2    string `form:"ticketType2" `
		SeatType2      string `form:"seatType2" `
		IDType2        string `form:"IDType2" `
		IDNumber2      string `form:"IDNumber2" `
	}
	P3 struct {
		PassengerName3 string `form:"passengerName3" `
		TicketType3    string `form:"ticketType3" `
		SeatType3      string `form:"seatType3" `
		IDType3        string `form:"IDType3" `
		IDNumber3      string `form:"IDNumber3" `
	}
	P4 struct {
		PassengerName4 string `form:"passengerName4" `
		TicketType4    string `form:"ticketType4" `
		SeatType4      string `form:"seatType4" `
		IDType4        string `form:"IDType4" `
		IDNumber4      string `form:"IDNumber4" `
	}
	P5 struct {
		PassengerName5 string `form:"passengerName5"`
		TicketType5    string `form:"ticketType5"`
		SeatType5      string `form:"seatType5"`
		IDType5        string `form:"IDType5"`
		IDNumber5      string `form:"IDNumber5"`
	}
}

func main() {
	m := martini.Classic()
	// render html templates from templates directory
	// m.Use(render.Renderer())
	m.Use(render.Renderer(render.Options{
		Directory: "templates", // Specify what path to load the templates from.
		// Layout:     "layout",          // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".html"}, // Specify extensions to load for templates.
		// Funcs:      []template.FuncMap{render.AppHelpers}, // Specify helper function maps for templates to access.
		Delims:     render.Delims{"{[{", "}]}"}, // Sets delimiters to the specified strings.
		Charset:    "UTF-8",                     // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                        // Output human readable JSON
	}))

	m.Get("/", func(r render.Render) {
		go LoginInit("")
		r.HTML(200, "login", nil)
	})

	m.Get("/loginPassCodeNew/**", loginPassCodeNewFunc)
	m.Get("/submitPassCodeNew/**", submitPassCodeNewFunc)
	m.Post("/login", binding.Form(UserLoginForm{}), LoginForm)
	m.Post("/query", binding.Form(TicketQuery{}), QueryForm)
	m.Post("/loadUser", loadUser)
	m.Get("/sock", Sock)
	// nodeWebkit, err := nw.New()
	// if err != nil {
	// 	panic(err)
	// }
	Info("a")
	// Pick a random localhost port, start listening for http requests using default handler
	// and send a message back to node-webkit to redirect
	// if err := nodeWebkit.ListenAndServe(m); err != nil {
	// 	panic(err)
	// }
	Info("b")
	m.Run()
	Info("c")
	// log.Fatal(http.ListenAndServe(":8080", m))
}
func Sock(w http.ResponseWriter, r *http.Request) {
	var err error
	ws, err = websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			log.Println("bye")
			log.Println(err)
			return
		}
		msg := string(p)
		Info(msg)
		if strings.Contains(msg, "code#") {
			code := msg[5:]
			Info("code:", code)

			if b, msg := order.checkRandCodeAnsyn(code); !b {
				Info(msg)
				ws.WriteMessage(1, []byte("update"))
				return
			}

			go func() {
				order.checkOrderInfo()
			}()

			order.SubmitCaptchaStr <- code
		}
		Info(messageType)
	}
}
func submitPassCodeNewFunc(res http.ResponseWriter, req *http.Request, params martini.Params) {

	req, err := http.NewRequest("GET", URLPassCodeNewPassenger+"&"+params["_1"], nil)
	if err != nil {
		Error("submitPassCodeNewFunc http.NewRequest error:", err)
		return
	}

	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/confirmPassenger/initDc"}
	AddReqestHeader(req, "GET", h)

	resp, err := client.Do(req)
	if err != nil {
		Error("submitPassCodeNewFunc client.Do error:", err)
		return
	}
	defer resp.Body.Close()

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Error("ioutil.ReadAll:", err)
		return
	}
	res.Header().Set("Content-Type", "image/jpeg")
	res.Write(bodyByte)
}
func loginPassCodeNewFunc(res http.ResponseWriter, req *http.Request, params martini.Params) {
	login.setNewCookie()

	req, err := http.NewRequest("GET", URLLoginPassCode+"&"+params["_1"], nil)
	if err != nil {
		Error("setNewCookie http.NewRequest error:", err)
		return
	}

	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/leftTicket/init"}
	AddReqestHeader(req, "GET", h)

	resp, err := client.Do(req)
	if err != nil {
		Error("setNewCookie client.Do error:", err)
		return
	}
	defer resp.Body.Close()

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Error("ioutil.ReadAll:", err)
		return
	}
	res.Header().Set("Content-Type", "image/jpeg")
	res.Write(bodyByte)
}
func QueryForm(res http.ResponseWriter, req *http.Request, params martini.Params, r render.Render, tq TicketQuery) {
	Info(tq)
	tq.parseTicket()
	tq.parseStranger()
	Info(tq)
	go func() {
		// LeftTicketInit()
		// DYQueryJs()

		order = tq.Order()
		Info(order)
		// q.leftTicketInit()
		checkUser("")
		order.submitOrderRequest()
		// q.DYQueryJs()
		order.initDc()
		order.GetPassengerDTO()
		Info("order:", order)
		// go func() {
		// 	order.checkOrderInfo()
		// }()
		ws.WriteMessage(1, []byte("update"))
	}()

	r.JSON(200, map[string]interface{}{"r": true, "o": tq})
}
func LoginForm(res http.ResponseWriter, req *http.Request, params martini.Params, r render.Render, l UserLoginForm) {
	fmt.Println(l.Username)
	fmt.Println(l.Password)
	fmt.Println(l.Code)
	login.Username = l.Username
	login.Password = l.Password
	login.Captcha = l.Code
	if b, err := login.checkUser(); err != nil {
		Error(err)
		return
	} else if b {
		Info("have logined!")
		r.HTML(200, "main", nil)
		return
	}

	if result, msg := login.CheckRandCodeAnsyn(); !result {
		r.HTML(200, "login", map[string]interface{}{"r": !result, "msg": msg, "username": l.Username, "password": l.Password})
		return
	}
	if result, msg := login.loginAysnSuggest(); !result {
		r.HTML(200, "login", map[string]interface{}{"r": !result, "msg": msg, "username": l.Username, "password": l.Password})
		return
	}

	login.checkUser()
	login.userLogin()
	go login.initQueryUserInfo()
	login.leftTicketInit()
	r.HTML(200, "main", nil)
}
func loadUser(res http.ResponseWriter, req *http.Request, params martini.Params, r render.Render) {
	if b, err := login.checkUser(); !b {
		r.JSON(200, map[string]interface{}{"r": b, "o": err})
		return
	}
	login.leftTicketInit()
	dyQueryJs("")
	GetPassCodes("")
	// login.userLogin()
	passenger := login.getPassengerDTO()

	r.JSON(200, map[string]interface{}{"r": true, "o": passenger.Data.NormalPassengers})
}
