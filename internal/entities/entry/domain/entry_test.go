package domain_test

import (
	"context"
	"strings"
	"testing"

	"github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/require"
)

func TestEntry_Validate(t *testing.T) {
	tests := []struct {
		name       string
		entry      *domain.Entry
		violations []validationtest.ViolationAttributes
	}{
		{
			name: "valid_entry",
			entry: &domain.Entry{
				Title:   "Valid Title",
				Content: "Valid Content",
			},
		},
		{
			name: "blank_title",
			entry: &domain.Entry{
				Title:   "",
				Content: "Content",
			},
			violations: []validationtest.ViolationAttributes{
				{Error: validation.ErrIsBlank, PropertyPath: "title"},
			},
		},
		{
			name: "title_too_long",
			entry: &domain.Entry{
				Title:   strings.Repeat("x", domain.MaxTitleLength+1),
				Content: "Content",
			},
			violations: []validationtest.ViolationAttributes{
				{Error: validation.ErrTooLong, PropertyPath: "title"},
			},
		},
		{
			name: "content_too_long",
			entry: &domain.Entry{
				Title:   "Title",
				Content: strings.Repeat("x", domain.MaxContentLength+1),
			},
			violations: []validationtest.ViolationAttributes{
				{Error: validation.ErrTooLong, PropertyPath: "content"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.entry.Validate(context.Background(), validator.Instance())

			if len(test.violations) == 0 {
				require.NoError(t, err)
			} else {
				validationtest.Assert(t, err).IsViolationList().WithAttributes(test.violations...)
			}
		})
	}
}
