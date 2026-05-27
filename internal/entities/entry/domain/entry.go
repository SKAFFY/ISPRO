package domain

import (
	"context"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

const MaxTitleLength = 255
const MaxContentLength = 10000

type Entry struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *Entry) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(ctx,
		validation.StringProperty("title", e.Title,
			it.IsNotBlank(),
			it.HasMaxLength(MaxTitleLength),
		),
		validation.StringProperty("content", e.Content,
			it.HasMaxLength(MaxContentLength),
		),
	)
}
