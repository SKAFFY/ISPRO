package factories

import (
	"context"
	"fmt"

	"github.com/go-openapi/loads"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	lookup "github.com/mkheyfets/ispro-app/internal/di/lookup"
	entry "github.com/mkheyfets/ispro-app/internal/entities/entry"
	entryRepo "github.com/mkheyfets/ispro-app/internal/entities/entry/repository"
	link "github.com/mkheyfets/ispro-app/internal/entities/link"
	linkRepo "github.com/mkheyfets/ispro-app/internal/entities/link/repository"
	"github.com/mkheyfets/ispro-app/internal/restapi"
	"github.com/mkheyfets/ispro-app/internal/restapi/operations"
	"github.com/mkheyfets/ispro-app/internal/restapi/operations/entries"
	"github.com/mkheyfets/ispro-app/internal/restapi/operations/links"
	"github.com/muonsoft/validation"
)

func CreatePool(ctx context.Context, c lookup.Container) (*pgxpool.Pool, error) {
	dsn := c.DSN(ctx)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return pool, nil
}

func CreateValidator(ctx context.Context, c lookup.Container) (*validation.Validator, error) {
	return validation.NewValidator()
}

func CreateEntryRepository(ctx context.Context, c lookup.Container) (*entryRepo.PostgresRepository, error) {
	pool := c.Pool(ctx)
	return entryRepo.NewPostgresRepository(pool), nil
}

func CreateLinkRepository(ctx context.Context, c lookup.Container) (*linkRepo.PostgresLinkRepository, error) {
	pool := c.Pool(ctx)
	return linkRepo.NewPostgresLinkRepository(pool), nil
}

func CreateEntryHandler(ctx context.Context, c lookup.Container) (*entry.Handler, error) {
	repo := c.EntryRepository(ctx)
	validator := c.Validator(ctx)

	listEntries := entry.NewListEntriesUseCase(repo)
	getEntry := entry.NewGetEntryUseCase(repo)
	createEntry := entry.NewCreateEntryUseCase(repo, validator)
	updateEntry := entry.NewUpdateEntryUseCase(repo, validator)
	deleteEntry := entry.NewDeleteEntryUseCase(repo)
	return entry.NewHandler(listEntries, getEntry, createEntry, updateEntry, deleteEntry), nil
}

func CreateLinkHandler(ctx context.Context, c lookup.Container) (*link.Handler, error) {
	linkRepository := c.LinkRepository(ctx)
	entryRepository := c.EntryRepository(ctx)
	validator := c.Validator(ctx)

	listLinks := link.NewListLinksUseCase(linkRepository)
	getLink := link.NewGetLinkUseCase(linkRepository)
	createLink := link.NewCreateLinkUseCase(linkRepository, validator, entryRepository)
	updateLink := link.NewUpdateLinkUseCase(linkRepository, validator, entryRepository)
	deleteLink := link.NewDeleteLinkUseCase(linkRepository)
	return link.NewHandler(listLinks, getLink, createLink, updateLink, deleteLink), nil
}

func CreateAPI(ctx context.Context, c lookup.Container) (*operations.IsproAppAPI, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load swagger spec: %w", err)
	}
	return operations.NewIsproAppAPI(swaggerSpec), nil
}

func CreateServer(ctx context.Context, c lookup.Container) (*restapi.Server, error) {
	api := c.API(ctx)
	entryHandler := c.EntryHandler(ctx)
	linkHandler := c.LinkHandler(ctx)

	api.EntriesListEntriesHandler = entries.ListEntriesHandlerFunc(entryHandler.ListEntries)
	api.EntriesCreateEntryHandler = entries.CreateEntryHandlerFunc(entryHandler.CreateEntry)
	api.EntriesGetEntryHandler = entries.GetEntryHandlerFunc(entryHandler.GetEntry)
	api.EntriesUpdateEntryHandler = entries.UpdateEntryHandlerFunc(entryHandler.UpdateEntry)
	api.EntriesDeleteEntryHandler = entries.DeleteEntryHandlerFunc(entryHandler.DeleteEntry)
	api.LinksListLinksHandler = links.ListLinksHandlerFunc(linkHandler.ListLinks)
	api.LinksCreateLinkHandler = links.CreateLinkHandlerFunc(linkHandler.CreateLink)
	api.LinksGetLinkHandler = links.GetLinkHandlerFunc(linkHandler.GetLink)
	api.LinksUpdateLinkHandler = links.UpdateLinkHandlerFunc(linkHandler.UpdateLink)
	api.LinksDeleteLinkHandler = links.DeleteLinkHandlerFunc(linkHandler.DeleteLink)

	return restapi.NewServer(api), nil
}

func CreateDSN(ctx context.Context, c lookup.Container) (string, error) {
	panic("not implemented")
}
