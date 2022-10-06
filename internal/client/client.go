package client

import (
	"context"
	"fmt"
	"github.com/trezorg/pow/internal/config"
	"github.com/trezorg/pow/internal/log"
	"github.com/trezorg/pow/internal/pow"
	"io"
	"net"
	"sync"
	"time"
)

type responseHandler func([]byte)

func DefaultHandler(response []byte) {
	log.Infof("Got response: %v", string(response))
}

func Request(address string, port int, zeroes int, handler responseHandler) error {
	serverAddr := fmt.Sprintf("%s:%d", address, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("resolveTCPAddr failed: %w", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Errorf("Error closing connection: %v", err)
		}

	}()

	puzzle := pow.Puzzle{}
	data, err := puzzle.MarshalBinary()
	if err != nil {
		return fmt.Errorf("error marshalling: %w", err)
	}
	log.Info("Sending request for puzzle...")
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("write to server failed: %w", err)
	}

	err = puzzle.Read(conn)
	if err != nil {
		return fmt.Errorf("read puzzle from server failed: %w", err)
	}

	responsePuzzle := pow.Find(puzzle, zeroes)
	data, err = responsePuzzle.MarshalBinary()
	if err != nil {
		return fmt.Errorf("error marshalling: %w", err)
	}

	log.Info("Sending solved puzzle...")
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("write to server failed: %w", err)
	}

	response, err := io.ReadAll(conn)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	handler(response)
	return nil
}

func Start(ctx context.Context, wg *sync.WaitGroup, cfg config.Config, handler responseHandler) {

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(cfg.Period):
			if err := Request(cfg.Address, cfg.Port, cfg.Zeroes, handler); err != nil {
				log.Error(err)
			}
		}
	}
}
