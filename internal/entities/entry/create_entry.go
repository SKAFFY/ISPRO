package entry

import (
	"context"

	"github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
	"github.com/muonsoft/validation"
)

type CreateEntryUseCase struct {
	repo      domain.EntryRepository
	validator *validation.Validator
}

func NewCreateEntryUseCase(repo domain.EntryRepository, validator *validation.Validator) *CreateEntryUseCase {
	return &CreateEntryUseCase{
		repo:      repo,
		validator: validator,
	}
}

type CreateEntryCommand struct {
	Title   string
	Content string
}

func (uc *CreateEntryUseCase) Handle(ctx context.Context, cmd CreateEntryCommand) (*domain.Entry, error) {
	entry := &domain.Entry{
		Title:   cmd.Title,
		Content: cmd.Content,
	}

	if err := entry.Validate(ctx, uc.validator); err != nil {
		return nil, err
	}

	return uc.repo.Create(ctx, cmd.Title, cmd.Content)
}
