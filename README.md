<a href="https://godoc.org/github.com/nulijiabei/freews"><img src="https://godoc.org/github.com/nulijiabei/freews?status.svg" alt="GoDoc"></a>

-------------

	快速将类和子方法转换成WebSocket接口 ...

-------------

	为什么要做这个项目:
	
		现在的 WebSocket 实现方式一般 ...
			http.Handle("/hello", websocket.Handler(HelloHandler))
			func HelloHandler(ws *websocket.Conn) {}
		这样一来如果想 ...
			func (this *Service) HelloHandler(ws *websocket.Conn) {} 
		是不可能的 ... 你可能说我可以用全局 ... 
		但是你可以看看更好的 freews ... 不但可以使用你自定义的类而且支持多个类 ...
			WSConn = websocket.Conn 的高可移植性 ... 
		把你之前的实现复制粘贴过来即可 ... websocket 怎么用这里怎么用 ...
		
-------------
	
	类名称及函数名称均会被转化为小写

-------------

	WSConn == websocket.Conn 这样做只是为了减少引用websocket包

-------------

	ws://127.0.0.1:8080/demo/hello
	ws://127.0.0.1:8080/hello/nihao

-------------

	package gcpool
	
	import (
		"log"
		"testing"
	)
	
	type Demo struct {
	}
	
	func (this *Demo) Hello(conn *WSConn) {
		// WSConn = websocket.Conn
	}
	
	type HELLO struct {
	}
	
	func (this *HELLO) Nihao(conn *WSConn) {
		conn.Write([]byte("Welcome !!!"))
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
	
	func Test(t *testing.T) {
	
		// New Service
		service := NewService()
		// 注册类到服务
		service.Register(new(Demo))
		service.Register(new(HELLO))
		...
		// 启动服务
		err := service.Start(":8080")
		if err != nil {
			log.Panic(err)
		}
	
	}

-------------