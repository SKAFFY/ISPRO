package validation

import (
	"context"

	entryDomain "github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
	"github.com/muonsoft/validation"
)

type EntryExistsConstraint struct {
	entryRepo entryDomain.EntryRepository
}

func EntryExists(entryRepo entryDomain.EntryRepository) *EntryExistsConstraint {
	return &EntryExistsConstraint{entryRepo: entryRepo}
}

func (c *EntryExistsConstraint) ValidateComparable(ctx context.Context, v *validation.Validator, value *int64) error {
	return c.validateEntryExists(ctx, v, value)
}

func (c *EntryExistsConstraint) ValidateNumber(ctx context.Context, v *validation.Validator, value *int64) error {
	return c.validateEntryExists(ctx, v, value)
}

func (c *EntryExistsConstraint) validateEntryExists(ctx context.Context, v *validation.Validator, value *int64) error {
	if value == nil || *value == 0 {
		return nil
	}

	entry, err := c.entryRepo.Get(ctx, *value)
	if err != nil {
		return err
	}
	if entry == nil {
		return v.BuildViolation(ctx, validation.ErrNotValid, "Entry with this ID does not exist.").Create()
	}

	return nil
}
