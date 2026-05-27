package entry

import (
	"context"

	"github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
)

type ListEntriesUseCase struct {
	repo domain.EntryRepository
}

func NewListEntriesUseCase(repo domain.EntryRepository) *ListEntriesUseCase {
	return &ListEntriesUseCase{repo: repo}
}

func (uc *ListEntriesUseCase) Handle(ctx context.Context) ([]*domain.Entry, error) {
	return uc.repo.List(ctx)
}
