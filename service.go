package freews

import (
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"

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
	// 遍历路径
	for path, r := range this.router {
		// 匹配路径
		if req.URL.Path == path {
			// 遍历域
			for a1, do := range this.domain {
				// 匹配域名
				if strings.ToLower(a1) == strings.ToLower(r.name) {
					// 遍历函数
					for a2, method := range do.methods {
						// 匹配函数名
						if strings.ToLower(a2) == strings.ToLower(r.method) {
							// 标记匹配
							status = true
							// WebSocket
							s := websocket.Server{Handshake: websocket.CheckOrigin}
							s.ServeWebSocket(w, req, func(conn *WSConn) {
								// 创建参数集
								value := make([]reflect.Value, method.Type.NumIn())
								// 第一个值为类反射值
								value[0] = do.rcvr
								// websocket.Conn ...
								value[1] = reflect.ValueOf(conn)
								// 调用函数 ... 传参 ... 并获取返回值 ...
								method.Func.Call(value)
							})
						}
					}
				}
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
