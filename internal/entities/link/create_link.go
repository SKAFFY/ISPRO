package link

import (
	"context"

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
		return nil, err
	}

	if cmd.SourceID == cmd.TargetID {
		return nil, validation.ErrIsEqual
	}

	entryExistsConstraint := linkValidation.EntryExists(uc.entryRepo)

	if err := uc.validator.Validate(ctx,
		validation.NumberProperty[int64]("source_id", link.SourceID, entryExistsConstraint),
		validation.NumberProperty[int64]("target_id", link.TargetID, entryExistsConstraint),
	); err != nil {
		return nil, err
	}

	result, err := uc.repo.Create(ctx, cmd.SourceID, cmd.TargetID)
	if err == nil {
		metrics.LinksCreatedTotal.Inc()
	}
	return result, err
}
