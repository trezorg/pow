package main

import (
	"context"
	"github.com/trezorg/pow/internal/client"
	"github.com/trezorg/pow/internal/config"
	"github.com/trezorg/pow/internal/log"
	"github.com/trezorg/pow/internal/server"
	"os"
	"os/signal"
	"sync"
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

	sig := <-stop
	log.Infof("Stopping server. Signal %v", sig)
	srv.Stop()
	close(stop)
}

func startClient(cfg config.Config) {
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go client.Start(ctx, &wg, cfg, client.DefaultHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop
	log.Infof("Stopping client. Signal %v", sig)
	cancel()
	wg.Wait()
	close(stop)
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
