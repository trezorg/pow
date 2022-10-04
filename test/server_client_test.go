package test

import (
	"github.com/stretchr/testify/require"
	"github.com/trezorg/pow/internal/client"
	"github.com/trezorg/pow/internal/log"
	"github.com/trezorg/pow/internal/server"
	"net"
	"testing"
	"time"
)

func freePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := l.Close(); err != nil {
			log.Error(err)
		}

	}()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func TestServerClient(t *testing.T) {
	log.NewLocalZapLogger("DEBUG")
	port, err := freePort()
	address := "127.0.0.1"
	require.NoError(t, err)
	srv, err := server.New(address, port, 5)
	require.NoError(t, err)
	defer srv.Stop()

	var response string

	handler := func(bs []byte) {
		response = string(bs)
	}

	go client.Start(address, port, 5, handler)

loop:
	for i := 0; i < 10; i++ {
		select {
		case <-time.After(time.Second):
			if response != "" {
				break loop
			}
		}
	}

	require.Greater(t, len(response), 0)

}
