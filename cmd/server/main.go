package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mkheyfets/ispro-app/internal/di"
	"github.com/mkheyfets/ispro-app/internal/logs"
	"github.com/mkheyfets/ispro-app/internal/tracing"
)

func main() {
	lokiEndpoint := os.Getenv("LOKI_ENDPOINT")
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "ispro-app"
	}

	logger := slog.New(logs.NewLokiHandler(lokiEndpoint, serviceName, slog.LevelInfo))
	slog.SetDefault(logger)

	ctx := context.Background()

	tp, err := tracing.InitTracerProvider(ctx, serviceName)
	if err != nil {
		slog.Error("failed to init tracer provider", "error", err)
		os.Exit(1)
	}
	defer tracing.ShutdownTracerProvider(ctx, tp)

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
		if err := server.Serve(); err != nil && err != http.ErrServerClosed {
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
