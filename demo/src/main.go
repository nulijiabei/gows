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
}

type HELLO struct {
}

func (this *HELLO) Nihao(conn *freews.WSConn) {
	// WSConn = websocket.Conn
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

func main() {

	service := freews.NewService()
	service.Register(new(Demo))
	service.Register(new(HELLO))
	err := service.Start(":8080")
	if err != nil {
		log.Panic(err)
	}

}
