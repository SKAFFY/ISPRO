package definitions

import (
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	entryHandler "github.com/mkheyfets/ispro-app/internal/entities/entry"
	entryRepo "github.com/mkheyfets/ispro-app/internal/entities/entry/repository"
	linkHandler "github.com/mkheyfets/ispro-app/internal/entities/link"
	linkRepo "github.com/mkheyfets/ispro-app/internal/entities/link/repository"
	"github.com/mkheyfets/ispro-app/internal/restapi"
	"github.com/mkheyfets/ispro-app/internal/restapi/operations"
	"github.com/muonsoft/validation"
)

type Container struct {
	// di:set
	DSN string

	// di:public
	Pool *pgxpool.Pool

	// di:public
	Validator *validation.Validator

	// di:public
	EntryRepository *entryRepo.PostgresRepository

	// di:public
	LinkRepository *linkRepo.PostgresLinkRepository

	// di:public
	EntryHandler *entryHandler.Handler

	// di:public
	LinkHandler *linkHandler.Handler

	// di:public
	API *operations.IsproAppAPI

	// di:public
	Server *restapi.Server
}
