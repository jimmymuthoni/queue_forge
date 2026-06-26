package main

import (
	"context"
	"log"
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
	server := apiserver.New(conf)
	server.Start(ctx)

	if err := server.Start(ctx); err != nil {
		return err
	}
	return nil
}