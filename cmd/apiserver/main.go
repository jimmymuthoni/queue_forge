package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jimmymuthoni/queue_forge/apiserver"
	"github.com/jimmymuthoni/queue_forge/config"
)

func main(){
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	conf, err := config.New()
	if err != nil {
		return err
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)
	server := apiserver.New(conf,logger)

	if err := server.Start(ctx); err != nil {
		return err
	}
	return nil
}