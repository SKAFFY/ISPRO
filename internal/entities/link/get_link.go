package link

import (
	"context"

	"github.com/mkheyfets/ispro-app/internal/entities/link/domain"
)

type GetLinkUseCase struct {
	repo domain.LinkRepository
}

func NewGetLinkUseCase(repo domain.LinkRepository) *GetLinkUseCase {
	return &GetLinkUseCase{repo: repo}
}

type GetLinkQuery struct {
	ID int64
}

func (uc *GetLinkUseCase) Handle(ctx context.Context, query GetLinkQuery) (*domain.Link, error) {
	return uc.repo.Get(ctx, query.ID)
}
