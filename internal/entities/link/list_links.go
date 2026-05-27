package link

import (
	"context"

	"github.com/mkheyfets/ispro-app/internal/entities/link/domain"
)

type ListLinksUseCase struct {
	repo domain.LinkRepository
}

func NewListLinksUseCase(repo domain.LinkRepository) *ListLinksUseCase {
	return &ListLinksUseCase{repo: repo}
}

func (uc *ListLinksUseCase) Handle(ctx context.Context) ([]*domain.Link, error) {
	return uc.repo.List(ctx)
}
