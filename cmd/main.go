package main

import (
	"github.com/trezorg/pow/internal/client"
	"github.com/trezorg/pow/internal/config"
	"github.com/trezorg/pow/internal/log"
	"github.com/trezorg/pow/internal/server"
	"os"
	"os/signal"
	"syscall"
)

var version = "0.0.1"

func startServer(cfg config.Config) {
	srv, err := server.New(cfg.Address, cfg.Port, cfg.Zeroes)
	if err != nil {
		log.Errorf("Cannot create srv: %v", err)
		os.Exit(1)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	for sig := range stop {
		log.Infof("Stopping server. Signal %v", sig)
		srv.Stop()
		close(stop)
	}
}

func startClient(cfg config.Config) {
	client.Start(cfg.Address, cfg.Port, cfg.Zeroes, client.DefaultHandler)
}

func main() {
	cfg := config.New(version)
	zap := log.NewLocalZapLogger(cfg.LogLevel)
	defer zap()
	if cfg.Command == config.Server {
		startServer(cfg)
	}
	if cfg.Command == config.Client {
		startClient(cfg)
	}
}
