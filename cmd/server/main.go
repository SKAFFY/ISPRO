package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-openapi/loads"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mkheyfets/ispro-app/pkg/handler"
	"github.com/mkheyfets/ispro-app/pkg/repository"
	"github.com/mkheyfets/ispro-app/pkg/service"
	"github.com/mkheyfets/ispro-app/restapi"
	"github.com/mkheyfets/ispro-app/restapi/operations"
	"github.com/mkheyfets/ispro-app/restapi/operations/entries"
	"github.com/mkheyfets/ispro-app/restapi/operations/links"
)

func main() {
	ctx := context.Background()

	connStr := "postgres://ispro:ispro@localhost:5432/ispro?sslmode=disable"
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	repo := repository.New(pool)
	svc := service.New(repo)
	h := handler.New(svc)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalf("failed to load swagger spec: %v", err)
	}

	api := operations.NewIsproAppAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	api.Logger = log.Printf

	api.EntriesListEntriesHandler = entries.ListEntriesHandlerFunc(h.ListEntries)
	api.EntriesCreateEntryHandler = entries.CreateEntryHandlerFunc(h.CreateEntry)
	api.EntriesGetEntryHandler = entries.GetEntryHandlerFunc(h.GetEntry)
	api.EntriesUpdateEntryHandler = entries.UpdateEntryHandlerFunc(h.UpdateEntry)
	api.EntriesDeleteEntryHandler = entries.DeleteEntryHandlerFunc(h.DeleteEntry)
	api.LinksListLinksHandler = links.ListLinksHandlerFunc(h.ListLinks)
	api.LinksCreateLinkHandler = links.CreateLinkHandlerFunc(h.CreateLink)
	api.LinksGetLinkHandler = links.GetLinkHandlerFunc(h.GetLink)
	api.LinksUpdateLinkHandler = links.UpdateLinkHandlerFunc(h.UpdateLink)
	api.LinksDeleteLinkHandler = links.DeleteLinkHandlerFunc(h.DeleteLink)

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