package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQualityValidation_ShouldReturnError_WhenQualityIsOutOfRange(t *testing.T) {
	tests := []struct {
		name    string
		quality int
	}{
		{"below_min", 0},
		{"above_max", 101},
		{"negative", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &QualityValidation{Quality: tt.quality}
			err := v.Validate()
			assert.ErrorIs(t, err, ErrInvalidQuality)
		})
	}
}

func TestQualityValidation_ShouldSucceed_WhenQualityIsInRange(t *testing.T) {
	tests := []struct {
		name    string
		quality int
	}{
		{"min", 1},
		{"max", 100},
		{"mid", 75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &QualityValidation{Quality: tt.quality}
			err := v.Validate()
			assert.NoError(t, err)
		})
	}
}
