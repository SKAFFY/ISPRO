package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mkheyfets/ispro-app/internal/di"
)

func main() {
	ctx := context.Background()

	container, err := di.NewContainer(
		di.SetDSN("postgres://ispro:ispro@localhost:5432/ispro?sslmode=disable"),
	)
	if err != nil {
		log.Fatalf("failed to create container: %v", err)
	}
	defer container.Close()

	server, err := container.Server(ctx)
	if err != nil {
		log.Fatalf("failed to get server: %v", err)
	}

	server.Host = "0.0.0.0"
	server.Port = 8080

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	time.Sleep(1 * time.Second)
}
