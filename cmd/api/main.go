package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/LeonardoForconesi/go-template/internal/config"
	"github.com/LeonardoForconesi/go-template/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	app, err := server.Build(cfg)
	if err != nil {
		log.Fatalf("build: %v", err)
	}

	go func() {
		if err := app.HTTP.ListenAndServe(); err != nil {
			log.Printf("http server stopped: %v", err)
		}
	}()

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	_ = app.HTTP.Shutdown(context.Background())
}
