package gows

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"runtime"

	"./websocket"
)

type WSConn = websocket.Conn

type Domain struct {
	name    string
	rcvr    reflect.Value
	typ     reflect.Type
	methods map[string]reflect.Method
}

type R struct {
	name   string
	method string
}

type Service struct {
	domain map[string]*Domain
	router map[string]R
}

// 创建 Service
func NewService() *Service {
	service := new(Service)
	service.domain = make(map[string]*Domain)
	service.router = make(map[string]R)
	return service
}

// 添加路由
func (this *Service) Router(path string, method interface{}) {
	name := runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name()
	a1 := regexp.MustCompile(`\(\*.*\)`).FindAllString(name, -1)[0]
	a2 := regexp.MustCompile(`\).*-`).FindAllString(name, -1)[0]
	this.router[path] = R{a1[2 : len(a1)-1], a2[2 : len(a2)-1]}
}

// 将类及方法注册到
func (this *Service) Register(rcvr interface{}) {
	// Domain 初始化
	do := new(Domain)
	do.methods = make(map[string]reflect.Method)
	// 获取类的反射类型
	do.typ = reflect.TypeOf(rcvr)
	// 获取类的反射值
	do.rcvr = reflect.ValueOf(rcvr)
	// 获取类名 ...
	do.name = do.rcvr.Elem().Type().Name()
	// 遍历函数列表 ...
	for m := 0; m < do.typ.NumMethod(); m++ {
		// 获取函数
		method := do.typ.Method(m)
		// mtype := method.Type
		mname := method.Name
		// 以函数名为KEY存储函数 ...
		do.methods[mname] = method
	}
	// Reg to Service
	this.domain[do.name] = do
}

func (this *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 是否匹配
	status := false
	// 匹配路径
	if r, ok := this.router[req.URL.Path]; ok {
		// 匹配域
		if do, ok := this.domain[r.name]; ok {
			// 匹配函数
			if method, ok := do.methods[r.method]; ok {
				// 标记匹配
				status = true
				// 创建参数集
				value := make([]reflect.Value, method.Type.NumIn())
				// 第一个值为类反射值
				value[0] = do.rcvr
				// WebSocket
				if method.Type.NumIn() > 1 {
					inType := method.Type.In(1).String()
					switch inType {
					case "*websocket.Conn":
						s := websocket.Server{Handshake: websocket.CheckOrigin}
						s.ServeWebSocket(w, req, func(conn *WSConn) {
							value[1] = reflect.ValueOf(conn)
							method.Func.Call(value)
						})
						return
					}
				}
				// 遍历注册函数参数类型
				for n := 1; n < method.Type.NumIn(); n++ {
					// 获取参数类型
					inType := method.Type.In(n).String()
					// 匹配参数类型
					switch inType {
					case "http.ResponseWriter":
						value[n] = reflect.ValueOf(w)
					case "*http.Request":
						value[n] = reflect.ValueOf(req)
					default:
						fmt.Printf("unsupported in type: %s\n", inType)
					}
				}
				// 调用函数 ... 传参 ...
				method.Func.Call(value)
			}
		}
	}
	// 不匹配则 NotFound ...
	if !status {
		http.NotFound(w, req)
	}
}

// 启动服务
func (this *Service) Start(port string) error {
	// ListenAndServe(addr string, handler Handler)
	return http.ListenAndServe(port, this)
}
