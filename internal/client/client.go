package client

import (
	"fmt"
	"github.com/trezorg/pow/internal/log"
	"github.com/trezorg/pow/internal/pow"
	"io"
	"net"
	"os"
)

func Start(address string, port int, zeroes int) {
	serverAddr := fmt.Sprintf("%s:%d", address, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		log.Errorf("ResolveTCPAddr failed: %v", err)
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Errorf("Dial failed: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Errorf("Error closing connection: %v", err)
		}

	}()

	puzzle := pow.Puzzle{}
	data, err := puzzle.Marshal()
	if err != nil {
		log.Errorf("Error marshalling: %v", err)
		os.Exit(1)
	}
	log.Info("Sending request for puzzle...")
	_, err = conn.Write(data)
	if err != nil {
		log.Errorf("Write to server failed: %v", err)
		os.Exit(1)
	}

	err = puzzle.Read(conn)
	if err != nil {
		log.Errorf("Read puzzle from server failed: %v", err)
		os.Exit(1)
	}

	responsePuzzle := pow.Find(puzzle, zeroes)
	data, err = responsePuzzle.Marshal()
	if err != nil {
		log.Errorf("Error marshalling: %v", err)
		os.Exit(1)
	}

	log.Info("Sending solved puzzle...")
	_, err = conn.Write(data)
	if err != nil {
		log.Errorf("Write to server failed: %v", err)
		os.Exit(1)
	}

	response, err := io.ReadAll(conn)
	if err != nil {
		log.Errorf("Error reading response: %v", err)
		os.Exit(1)
	}

	log.Infof("Got response: %v", string(response))
}
