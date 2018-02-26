package gcpool

import (
	"bufio"
	"io"
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

	service := NewService()
	service.Register(new(Demo))
	service.Register(new(HELLO))
	err := service.Start(":8080")
	if err != nil {
		log.Panic(err)
	}

}
