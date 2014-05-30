package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

var (
	client = &http.Client{}
)

func init() {
	pr, err := url.Parse(Conf.ProxyUrl)
	if err != nil {
		Error(err)
		return
	}
	client = &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(pr),
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
		},
	}

}

type UserLoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Code     string `form:"code" binding:"required"`
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
		r.HTML(200, "login", nil)
	})

	m.Get("/loginPassCodeNew/**", func(res http.ResponseWriter, req *http.Request, params martini.Params) {

		rsp, err := client.Get(URLLoginPassCode + "&" + params["_1"])
		if err != nil {
			fmt.Println("client.Get:", err)
			return
		}
		defer rsp.Body.Close()

		bodyByte, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			Error("ioutil.ReadAll:", err)
			return
		}
		res.Header().Set("Content-Type", "image/jpeg")
		res.Write(bodyByte)
	})
	m.Post("/login", binding.Form(UserLoginForm{}), LoginForm)
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

func LoginForm(res http.ResponseWriter, req *http.Request, params martini.Params, r render.Render, l UserLoginForm) {
	fmt.Println(l.Username)
	fmt.Println(l.Password)
	fmt.Println(l.Code)
	r.HTML(200, "main", nil)
}
