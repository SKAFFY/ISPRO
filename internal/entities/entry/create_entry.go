package entry

import (
	"context"
	"log/slog"

	"github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
	"github.com/mkheyfets/ispro-app/internal/metrics"
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
		slog.Warn("entry validation failed", "title", cmd.Title)
		return nil, err
	}

	result, err := uc.repo.Create(ctx, cmd.Title, cmd.Content)
	if err != nil {
		slog.Error("failed to create entry in repository", "error", err, "title", cmd.Title)
		return nil, err
	}

	metrics.EntriesCreatedTotal.Inc()
	slog.Info("entry created", "id", result.ID, "title", result.Title)
	return result, nil
}
