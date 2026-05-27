package link

import (
	"context"
	"log/slog"

	entryDomain "github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
	"github.com/mkheyfets/ispro-app/internal/entities/link/domain"
	linkValidation "github.com/mkheyfets/ispro-app/internal/entities/link/domain/validation"
	"github.com/mkheyfets/ispro-app/internal/metrics"
	"github.com/muonsoft/validation"
)

type CreateLinkUseCase struct {
	repo      domain.LinkRepository
	validator *validation.Validator
	entryRepo entryDomain.EntryRepository
}

func NewCreateLinkUseCase(repo domain.LinkRepository, validator *validation.Validator, entryRepo entryDomain.EntryRepository) *CreateLinkUseCase {
	return &CreateLinkUseCase{
		repo:      repo,
		validator: validator,
		entryRepo: entryRepo,
	}
}

type CreateLinkCommand struct {
	SourceID int64
	TargetID int64
}

func (uc *CreateLinkUseCase) Handle(ctx context.Context, cmd CreateLinkCommand) (*domain.Link, error) {
	link := &domain.Link{
		SourceID: cmd.SourceID,
		TargetID: cmd.TargetID,
	}

	if err := link.Validate(ctx, uc.validator); err != nil {
		slog.Warn("link validation failed", "source_id", cmd.SourceID, "target_id", cmd.TargetID)
		return nil, err
	}

	if cmd.SourceID == cmd.TargetID {
		slog.Warn("link source equals target", "source_id", cmd.SourceID)
		return nil, validation.ErrIsEqual
	}

	entryExistsConstraint := linkValidation.EntryExists(uc.entryRepo)

	if err := uc.validator.Validate(ctx,
		validation.NumberProperty[int64]("source_id", link.SourceID, entryExistsConstraint),
		validation.NumberProperty[int64]("target_id", link.TargetID, entryExistsConstraint),
	); err != nil {
		slog.Warn("entry existence check failed", "source_id", cmd.SourceID, "target_id", cmd.TargetID)
		return nil, err
	}

	result, err := uc.repo.Create(ctx, cmd.SourceID, cmd.TargetID)
	if err != nil {
		slog.Error("failed to create link in repository", "error", err, "source_id", cmd.SourceID, "target_id", cmd.TargetID)
		return nil, err
	}

	metrics.LinksCreatedTotal.Inc()
	slog.Info("link created", "id", result.ID, "source_id", result.SourceID, "target_id", result.TargetID)
	return result, nil
}
