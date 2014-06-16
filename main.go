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
	"sync"
	"time"

	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

var (
	queryChannel = make(chan int, 3)    // 查询线程
	cdnChannel   = make(chan string, 1) // CDN线程
	querywg      = sync.WaitGroup{}     // 用于等待所有 goroutine 结束
	login        = &Login{}
	order        = &Order{}
	client       = &http.Client{}
	ws           *websocket.Conn
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
		go func() {
			loginInit(Conf.CDN[0])
			dynamicJsLoginJs(Conf.CDN[0])
		}()
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
	// Info("a")
	// Pick a random localhost port, start listening for http requests using default handler
	// and send a message back to node-webkit to redirect
	// if err := nodeWebkit.ListenAndServe(m); err != nil {
	// 	panic(err)
	// }
	Info("b")
	m.Run()
	Info("c")
}

// sock 接口
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
		Info("websocket:", msg, messageType)
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
				order.getQueueCount()
				order.confirmSingleForQueue()
			}()

			order.SubmitCaptchaStr <- code
		}
	}
}

//提交订单验证码
func submitPassCodeNewFunc(res http.ResponseWriter, req *http.Request, params martini.Params) {

	bodyByte := make([]byte, 30)
	defer func() {
		res.Header().Set("Content-Type", "image/jpeg")
		res.Write(bodyByte)
	}()

	req, err := http.NewRequest("GET", URLPassCodeNewPassenger+"&"+params["_1"], nil)
	if err != nil {
		Error("submitPassCodeNewFunc http.NewRequest error:", err)
		return
	}

	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/confirmPassenger/initDc"}
	AddReqestHeader(req, "GET", h)

	con, err := NewForwardClientConn(order.CDN, req.URL.Scheme)
	if err != nil {
		Error("DoForWardRequestHeader NewForwardClientConn error:", err)
		return
	}
	defer con.Close()
	resp, err := con.Do(req)

	if err != nil {
		Error("submitPassCodeNewFunc client.Do error:", err)
		return
	}
	defer resp.Body.Close()

	var err1 error
	bodyByte, err1 = ioutil.ReadAll(resp.Body)
	if err1 != nil {
		Error("ioutil.ReadAll:", err1)
		return
	}
}

//登陆验证码
func loginPassCodeNewFunc(res http.ResponseWriter, req *http.Request, params martini.Params) {
	login.setNewCookie()

	bodyByte := make([]byte, 30)
	defer func() {
		res.Header().Set("Content-Type", "image/jpeg")
		res.Write(bodyByte)
	}()

	req, err := http.NewRequest("GET", URLLoginPassCode+"&"+params["_1"], nil)
	if err != nil {
		Error("loginPassCodeNewFunc http.NewRequest error:", err)
		return
	}

	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/leftTicket/init"}
	AddReqestHeader(req, "GET", h)

	con, err := NewForwardClientConn(Conf.CDN[0], req.URL.Scheme)
	if err != nil {
		Error("loginPassCodeNewFunc NewForwardClientConn error:", err)
		return
	}
	defer con.Close()
	resp, err := con.Do(req)

	// resp, err := client.Do(req)
	if err != nil {
		Error("loginPassCodeNewFunc error:", err)
		return
	}
	defer resp.Body.Close()

	var err1 error
	bodyByte, err1 = ioutil.ReadAll(resp.Body)
	if err1 != nil {
		Error("ioutil.ReadAll:", err1)
		return
	}
}

//查询按钮
func QueryForm(res http.ResponseWriter, req *http.Request, params martini.Params, r render.Render, tq TicketQuery) {
	Info(tq)
	query := &Query{}
	tq.parseTicket(query)
	tq.parseStranger(query)
	Info(query)
	go func() {
		close(cdnChannel)
		cdnChannel = make(chan string, 1)
		for {
			for _, cdn := range Conf.CDN {
				cdnChannel <- cdn
			}
		}
	}()
	go func() {

		for {
			if cdn, ok := <-cdnChannel; ok {
				querywg.Add(1)
				queryChannel <- 1
				go func() {
					defer func() {
						<-queryChannel
						querywg.Done()
					}()
					query.CDN = cdn
					order = query.Order()
					if order != nil {
						if re, err := order.checkUser(); re {
							order.submitOrderRequest()
							order.initDc()
							go order.GetPassengerDTO()
							ws.WriteMessage(1, []byte("update"))
						} else {
							Error("checkUser 失败!", err)
						}
					}
				}()

			} else {
				break
			}
		}
		querywg.Wait()
	}()

	r.JSON(200, map[string]interface{}{"r": true, "o": query})
}

//登陆
func LoginForm(res http.ResponseWriter, req *http.Request, params martini.Params, r render.Render, l UserLoginForm) {
	fmt.Println(l.Username)
	fmt.Println(l.Password)
	fmt.Println(l.Code)
	login.Username = l.Username
	login.Password = l.Password
	login.Captcha = l.Code
	login.CDN = Conf.CDN[0]
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
	go func() {
		login.userLogin()
		login.initQueryUserInfo()
		login.leftTicketInit()
		dynamicJsQueryJs(login.CDN)
		getPassCodeNewInQueryPage(login.CDN)
	}()
	r.HTML(200, "main", nil)
}

//获取用户
func loadUser(res http.ResponseWriter, req *http.Request, params martini.Params, r render.Render) {
	if b, err := login.checkUser(); !b {
		r.JSON(200, map[string]interface{}{"r": b, "o": err})
		return
	}
	passenger := login.getPassengerDTO()

	r.JSON(200, map[string]interface{}{"r": true, "o": passenger.Data.NormalPassengers})
}
