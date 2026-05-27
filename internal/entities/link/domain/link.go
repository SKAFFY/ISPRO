package domain

import (
	"context"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

type Link struct {
	ID        int64     `json:"id"`
	SourceID  int64     `json:"source_id"`
	TargetID  int64     `json:"target_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (l *Link) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(ctx,
		validation.NumberProperty[int64]("source_id", l.SourceID, it.IsPositive[int64]()),
		validation.NumberProperty[int64]("target_id", l.TargetID, it.IsPositive[int64]()),
	)
}
