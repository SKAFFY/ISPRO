package entry

import (
	"context"

	"github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
)

type DeleteEntryUseCase struct {
	repo domain.EntryRepository
}

func NewDeleteEntryUseCase(repo domain.EntryRepository) *DeleteEntryUseCase {
	return &DeleteEntryUseCase{repo: repo}
}

type DeleteEntryCommand struct {
	ID int64
}

func (uc *DeleteEntryUseCase) Handle(ctx context.Context, cmd DeleteEntryCommand) error {
	return uc.repo.Delete(ctx, cmd.ID)
}
