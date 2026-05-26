package domain_test

import (
	"context"
	"testing"

	"github.com/mkheyfets/ispro-app/internal/entities/link/domain"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/require"
)

func TestLink_Validate(t *testing.T) {
	tests := []struct {
		name       string
		link       *domain.Link
		violations []validationtest.ViolationAttributes
	}{
		{
			name: "valid_link",
			link: &domain.Link{
				SourceID: 1,
				TargetID: 2,
			},
		},
		{
			name: "zero_source_id",
			link: &domain.Link{
				SourceID: 0,
				TargetID: 2,
			},
			violations: []validationtest.ViolationAttributes{
				{Error: validation.ErrNotPositive, PropertyPath: "source_id"},
			},
		},
		{
			name: "negative_source_id",
			link: &domain.Link{
				SourceID: -1,
				TargetID: 2,
			},
			violations: []validationtest.ViolationAttributes{
				{Error: validation.ErrNotPositive, PropertyPath: "source_id"},
			},
		},
		{
			name: "zero_target_id",
			link: &domain.Link{
				SourceID: 1,
				TargetID: 0,
			},
			violations: []validationtest.ViolationAttributes{
				{Error: validation.ErrNotPositive, PropertyPath: "target_id"},
			},
		},
		{
			name: "both_zero",
			link: &domain.Link{
				SourceID: 0,
				TargetID: 0,
			},
			violations: []validationtest.ViolationAttributes{
				{Error: validation.ErrNotPositive, PropertyPath: "source_id"},
				{Error: validation.ErrNotPositive, PropertyPath: "target_id"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.link.Validate(context.Background(), validator.Instance())

			if len(test.violations) == 0 {
				require.NoError(t, err)
			} else {
				validationtest.Assert(t, err).IsViolationList().WithAttributes(test.violations...)
			}
		})
	}
}
