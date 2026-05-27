package link

import (
	"context"

	entryDomain "github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
	"github.com/mkheyfets/ispro-app/internal/entities/link/domain"
	linkValidation "github.com/mkheyfets/ispro-app/internal/entities/link/domain/validation"
	"github.com/muonsoft/validation"
)

type UpdateLinkUseCase struct {
	repo      domain.LinkRepository
	validator *validation.Validator
	entryRepo entryDomain.EntryRepository
}

func NewUpdateLinkUseCase(repo domain.LinkRepository, validator *validation.Validator, entryRepo entryDomain.EntryRepository) *UpdateLinkUseCase {
	return &UpdateLinkUseCase{
		repo:      repo,
		validator: validator,
		entryRepo: entryRepo,
	}
}

type UpdateLinkCommand struct {
	ID       int64
	SourceID int64
	TargetID int64
}

func (uc *UpdateLinkUseCase) Handle(ctx context.Context, cmd UpdateLinkCommand) (*domain.Link, error) {
	link := &domain.Link{
		ID:       cmd.ID,
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

	return uc.repo.Update(ctx, cmd.ID, cmd.SourceID, cmd.TargetID)
}
