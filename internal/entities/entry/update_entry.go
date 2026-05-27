package entry

import (
	"context"

	"github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
	"github.com/muonsoft/validation"
)

type UpdateEntryUseCase struct {
	repo      domain.EntryRepository
	validator *validation.Validator
}

func NewUpdateEntryUseCase(repo domain.EntryRepository, validator *validation.Validator) *UpdateEntryUseCase {
	return &UpdateEntryUseCase{
		repo:      repo,
		validator: validator,
	}
}

type UpdateEntryCommand struct {
	ID      int64
	Title   string
	Content string
}

func (uc *UpdateEntryUseCase) Handle(ctx context.Context, cmd UpdateEntryCommand) (*domain.Entry, error) {
	entry := &domain.Entry{
		ID:      cmd.ID,
		Title:   cmd.Title,
		Content: cmd.Content,
	}

	if err := entry.Validate(ctx, uc.validator); err != nil {
		return nil, err
	}

	return uc.repo.Update(ctx, cmd.ID, cmd.Title, cmd.Content)
}
