package main

import (
	"crypto/tls"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/nfnt/resize"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			// Proxy:           http.ProxyURL(pr),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
)

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		// conn, err := net.DialTimeout(netw, addr, cTimeout)
		fmt.Println(netw)
		fmt.Println(addr)
		conn, err := tls.Dial(netw, addr, &tls.Config{
			InsecureSkipVerify: true,
		})
		if err != nil {
			fmt.Println("ccc", err)
			return nil, err
		}
		//conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func NewTimeoutClient(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Client {

	return &http.Client{
		Transport: &http.Transport{
			//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Dial: TimeoutDialer(connectTimeout, readWriteTimeout),
		},
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
		r.HTML(200, "login", "jeremy")
	})

	m.Get("/loginPassCodeNew1", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
		res.Header().Add("key", "value")
		res.Header().Set("Content-Type", "image/jpeg")
		res.WriteHeader(200) // HTTP 200
	})
	m.Get("/loginPassCodeNew/**", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		res.Header().Set("Content-Type", "image/jpeg")
		// client = NewTimeoutClient(20*time.Second, 20*time.Second)
		c := http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				Dial: func(netw, addr string) (net.Conn, error) {
					deadline := time.Now().Add(10 * time.Second)
					c, err := net.DialTimeout(netw, Conf.CDN[0]+":443", time.Second*10)
					if err != nil {
						return nil, err
					}
					c.SetDeadline(deadline)
					return c, nil
				},
			},
		}
		rsp, err := c.Get(URLLoginPassCode + "&" + params["_1"])
		if err != nil {
			fmt.Println("aaa", err)
			return
		}
		defer rsp.Body.Close()

		bodyByte, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			Error("ioutil.ReadAll:", err)
		}

		res.Write(bodyByte)
		res.WriteHeader(200)

	})
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
func thumb() image.Image {
	file, err := os.Open("test.jpg")
	if err != nil {
		log.Fatal(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	m := resize.Resize(0, 200, img, resize.MitchellNetravali)

	return m
}
