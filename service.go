package gcpool

import (
	"fmt"
	"net/http"
	"reflect"
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

// Server Json HTTP
type Service struct {
	domain map[string]*Domain
}

// 创建 Server 其中 def 为路径处理函数 nil 则使用默认
func NewService() *Service {
	service := new(Service)
	service.domain = make(map[string]*Domain)
	return service
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

func (this *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 是否匹配 ...
	status := false
	// 遍历域
	for a, do := range this.domain {
		// 域入口
		if strings.HasPrefix(strings.ToLower(r.URL.Path), fmt.Sprintf("/%s/", strings.ToLower(a))) {
			// 遍历函数
			for b, method := range do.methods {
				// 函数入口
				if strings.HasPrefix(strings.ToLower(r.URL.Path), fmt.Sprintf("/%s/%s", strings.ToLower(a), strings.ToLower(b))) {
					// 标记匹配
					status = true
					// WebSocket
					s := websocket.Server{Handshake: websocket.CheckOrigin}
					s.ServeWebSocket(w, r, func(conn *WSConn) {
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
	// 不匹配则 NotFound ...
	if !status {
		http.NotFound(w, r)
	}
}

// 启动服务
func (this *Service) Start(port string) error {
	// ListenAndServe(addr string, handler Handler)
	return http.ListenAndServe(port, this)
}
