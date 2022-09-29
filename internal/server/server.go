package server

import (
	"fmt"
	"github.com/trezorg/pow/internal/log"
	"net"
)

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
}

func Start(address string, port int) error {
	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Error(err)
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err)
			return err
		}
		go handleConnection(conn)
	}
}
