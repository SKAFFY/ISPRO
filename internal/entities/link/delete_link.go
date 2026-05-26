package link

import (
	"context"
	"errors"

	"github.com/mkheyfets/ispro-app/internal/entities/link/domain"
	"github.com/mkheyfets/ispro-app/internal/metrics"
)

var ErrLinkNotFound = errors.New("link not found")

type DeleteLinkUseCase struct {
	repo domain.LinkRepository
}

func NewDeleteLinkUseCase(repo domain.LinkRepository) *DeleteLinkUseCase {
	return &DeleteLinkUseCase{repo: repo}
}

type DeleteLinkCommand struct {
	ID int64
}

func (uc *DeleteLinkUseCase) Handle(ctx context.Context, cmd DeleteLinkCommand) error {
	existing, err := uc.repo.Get(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrLinkNotFound
	}

	err = uc.repo.Delete(ctx, cmd.ID)
	if err == nil {
		metrics.LinksDeletedTotal.Inc()
	}
	return err
}
