package main

import (
	"bufio"
	"io"
	"log"
	"net/http"

	"../../../gows"
)

type Demo struct {
}

func (this *Demo) Hello(conn *gows.WSConn) {
	// WSConn = websocket.Conn
	conn.Write([]byte("Baidu !!!"))
}

type HELLO struct {
}

func (this *HELLO) Nihao(conn *gows.WSConn) {
	// WSConn = websocket.Conn
	conn.Write([]byte("Sina !!!"))
	r := bufio.NewReader(conn)
	for {
		v, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		conn.Write(v)
		conn.Write([]byte("\n"))
	}
}

func (this *HELLO) Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Version 1.0"))
}

func main() {

	// New Service
	service := gows.NewService()
	// New Class
	demo := new(Demo)
	hello := new(HELLO)
	// 注册到服务
	service.Register(demo)
	service.Register(hello)
	// 添加路由
	service.Router("/v1/baidu", demo.Hello)
	service.Router("/v2/sina", hello.Nihao)
	service.Router("/api/version", hello.Version)
	// 启动服务
	err := service.Start(":8080")
	if err != nil {
		log.Panic(err)
	}

}
