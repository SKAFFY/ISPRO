package entry

import (
	"context"

	"github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
)

type GetEntryUseCase struct {
	repo domain.EntryRepository
}

func NewGetEntryUseCase(repo domain.EntryRepository) *GetEntryUseCase {
	return &GetEntryUseCase{repo: repo}
}

type GetEntryQuery struct {
	ID int64
}

func (uc *GetEntryUseCase) Handle(ctx context.Context, query GetEntryQuery) (*domain.Entry, error) {
	return uc.repo.Get(ctx, query.ID)
}
