package server

import (
	"fmt"
	"github.com/trezorg/pow/internal/log"
	"github.com/trezorg/pow/internal/pow"
	"github.com/trezorg/pow/internal/wordofwisdom"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PuzzleServer struct {
	listener      net.Listener
	quit          chan struct{}
	wg            sync.WaitGroup
	hashCashZeros int
}

func (server *PuzzleServer) handleConnection(conn net.Conn) {
	defer func() {
		log.Debugf("Closing connection from %s...", conn.RemoteAddr().String())
		if err := conn.Close(); err != nil {
			log.Error("Failed to close connection")
		}
	}()
	parts := strings.Split(conn.RemoteAddr().String(), ":")
	port, _ := strconv.ParseInt(parts[1], 10, 64)
	log.Debugf("Handling new connection from %s...", conn.RemoteAddr().String())

	var seed uint64

loop:
	for {
		select {
		case <-server.quit:
			return
		default:
			if err := conn.SetDeadline(time.Now().Add(time.Second)); err != nil {
				log.Errorf("Failed to set deadline: %v", err)
				return
			}
			puzzle := pow.Puzzle{}
			err := puzzle.Read(conn)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue loop
				} else if err != io.EOF {
					log.Errorf("Failed to read from connection: %v", err)
					return
				}
			}
			if puzzle.SeedID == 0 {
				// request for new puzzle
				req := pow.Request{
					RemoteIP:   parts[0],
					RemotePort: uint16(port),
					SeedID:     rand.Uint64(),
				}
				respPuzzle := pow.Make(req)
				resp, _ := respPuzzle.Marshal()
				_, err = conn.Write(resp)
				if err != nil {
					log.Errorf("Failed to write into connection: %v", err)
					return
				}
				seed = req.SeedID
			} else if seed == puzzle.SeedID {
				// solved puzzle
				req := pow.Request{
					RemoteIP:   parts[0],
					RemotePort: uint16(port),
					SeedID:     seed,
				}
				var message []byte
				if pow.Check(req, puzzle, server.hashCashZeros) {
					message = []byte(wordofwisdom.Quote())
				} else {
					message = []byte("fail")
				}
				_, err = conn.Write(message)
				if err != nil {
					log.Errorf("Failed to write into connection: %v", err)
				}
				return
			} else {
				log.Error("Inconsistent seed")
				return
			}
		}
	}
}

func New(address string, port int, hashCashZeros int) (*PuzzleServer, error) {
	adPort := fmt.Sprintf("%s:%d", address, port)
	listener, err := net.Listen("tcp", adPort)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	srv := &PuzzleServer{
		listener:      listener,
		quit:          make(chan struct{}),
		wg:            sync.WaitGroup{},
		hashCashZeros: hashCashZeros,
	}
	log.Infof("Starting tcp pow server %s...", adPort)
	srv.wg.Add(1)
	go srv.start()
	return srv, nil
}

func (server *PuzzleServer) start() {
	defer server.wg.Done()
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			select {
			case <-server.quit:
				log.Info("Exiting accept loop...")
				return
			default:
				log.Error(err)
				return
			}
		}
		server.wg.Add(1)
		go func() {
			server.handleConnection(conn)
			server.wg.Done()
		}()
	}
}

func (server *PuzzleServer) Stop() {
	close(server.quit)
	log.Info("Closing network listener...")
	if err := server.listener.Close(); err != nil {
		log.Error(err)
	}
	log.Debug("Waiting connections are being closed...")
	server.wg.Wait()
}
