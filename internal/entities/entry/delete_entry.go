package entry

import (
	"context"

	"github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
	"github.com/mkheyfets/ispro-app/internal/metrics"
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
	err := uc.repo.Delete(ctx, cmd.ID)
	if err == nil {
		metrics.EntriesDeletedTotal.Inc()
	}
	return err
}
