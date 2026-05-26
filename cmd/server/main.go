package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mkheyfets/ispro-app/internal/di"
	"github.com/mkheyfets/ispro-app/internal/logs"
)

func main() {
	vlEndpoint := os.Getenv("VICTORIALOGS_ENDPOINT")
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "ispro-app"
	}

	logger := slog.New(logs.NewVictoriaLogsHandler(vlEndpoint, serviceName, slog.LevelInfo))
	slog.SetDefault(logger)

	ctx := context.Background()

	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "postgres://ispro:ispro@localhost:5432/ispro?sslmode=disable"
	}

	container, err := di.NewContainer(
		di.SetDSN(dsn),
	)
	if err != nil {
		slog.Error("failed to create container", "error", err)
		os.Exit(1)
	}
	defer container.Close()

	server, err := container.Server(ctx)
	if err != nil {
		slog.Error("failed to get server", "error", err)
		os.Exit(1)
	}

	server.Host = "0.0.0.0"
	server.Port = 8080

	go func() {
		slog.Info("starting server", "host", server.Host, "port", server.Port)
		if err := server.Serve(); err != nil {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	time.Sleep(1 * time.Second)
}
