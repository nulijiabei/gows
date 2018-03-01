<a href="https://godoc.org/github.com/nulijiabei/freews"><img src="https://godoc.org/github.com/nulijiabei/freews?status.svg" alt="GoDoc"></a>

-------------

	快速将WebSocket集成到自定义类和子方法中 ...

-------------

	为什么要做这个项目:
	
		现在的 WebSocket 实现方式一般 ...
			http.Handle("/hello", websocket.Handler(HelloHandler))
			func HelloHandler(ws *websocket.Conn) {}
		这样一来如果想 ...
			func (this *MyClass) HelloHandler(ws *websocket.Conn) {
				this.MyData ... 增删改查 ... 等等 ...
			} 
		是不可能的 ... 你可能说我可以用全局 ... 
		但是你可以看看更好的 freews ... 不但可以使用你自定义的类而且支持多个类 ...
			WSConn = websocket.Conn 的高可移植性(只为了让你少引用websocket包) ... 
		把你之前的实现复制粘贴过来即可 ... websocket 怎么用这里怎么用 ...

-------------

	WSConn == websocket.Conn 这样做只是为了减少引用websocket包

-------------

	// 自定义访问地址 ...
	ws://127.0.0.1:8080/demo/hello
	ws://127.0.0.1:8080/hello/nihao
	ws://127.0.0.1:8080/v1/baidu
	ws://127.0.0.1:8080/v2/sina

-------------

	package main
	
	import (
		"bufio"
		"io"
		"log"
	
		"../../../freews"
	)
	
	type Demo struct {
	}
	
	func (this *Demo) Hello(conn *freews.WSConn) {
		// WSConn = websocket.Conn
		conn.Write([]byte("Baidu !!!"))
	}
	
	type HELLO struct {
	}
	
	func (this *HELLO) Nihao(conn *freews.WSConn) {
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
		service := freews.NewService()
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