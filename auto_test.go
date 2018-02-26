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
	// WSConn = websocket.Conn
}

func Test(t *testing.T) {

	service := NewService()
	service.Register(new(Demo))
	service.Register(new(HELLO))
	err := service.Start(":8080")
	if err != nil {
		log.Panic(err)
	}

}
