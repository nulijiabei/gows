<a href="https://godoc.org/github.com/nulijiabei/gows"><img src="https://godoc.org/github.com/nulijiabei/gows?status.svg" alt="GoDoc"></a>

-------------

	快速将自定义类及子方法转换为WebSocket-API
	
-------------

	为什么要做这个项目:
	
		现在的 WebSocket 实现方式一般 ...
			http.Handle("/hello", websocket.Handler(HelloHandler))
			func HelloHandler(ws *websocket.Conn) {}
		这样一来如果想 ...
			func (this *MyClass) HelloHandler(ws *websocket.Conn) {
				this.MyData ... 增删改查 ... 等等 ...
			}  
		是不可能的 ... 可能你还会有办法比如这样：... 等等 ...
			var MYCLASS *MyClass
			func HelloHandler(ws *websocket.Conn) {
				MYCLASS.MyData ... 增删改查 ... 等等 ...
			}
			... 但是 ...
		来看看 gows 吧 ...

-------------

	WSConn == websocket.Conn 这样做只是为了减少引用 websocket .
	
-------------

	// 自定义路由 ...
	ws://127.0.0.1:8080/v1/baidu
	ws://127.0.0.1:8080/v2/sina
	...

-------------

	package main
	
	import (
		"bufio"
		"io"
		"log"
	
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
		// 启动服务
		err := service.Start(":8080")
		if err != nil {
			log.Panic(err)
		}
	
	}

-------------