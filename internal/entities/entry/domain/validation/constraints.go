package validation

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"

	"github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
)

const MaxTitleLength = 255
const MaxContentLength = 10000

type EntryConstraints struct{}

func NewEntryConstraints() *EntryConstraints {
	return &EntryConstraints{}
}

func (c *EntryConstraints) Validate(ctx context.Context, validator *validation.Validator, entry *domain.Entry) error {
	return validator.Validate(ctx,
		validation.StringProperty("title", entry.Title,
			it.IsNotBlank(),
			it.HasMaxLength(MaxTitleLength),
		),
		validation.StringProperty("content", entry.Content,
			it.HasMaxLength(MaxContentLength),
		),
	)
}

func ValidateEntry(ctx context.Context, validator *validation.Validator, entry *domain.Entry) error {
	constraints := NewEntryConstraints()
	return constraints.Validate(ctx, validator, entry)
}
