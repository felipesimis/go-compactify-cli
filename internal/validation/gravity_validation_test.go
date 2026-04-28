package validation

import (
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/stretchr/testify/assert"
)

func TestGravityValidation_Validate(t *testing.T) {
	tests := []struct {
		name    string
		gravity image.Gravity
		wantErr error
	}{
		{"Should succeed at lower bound", image.GravityCentre, nil},
		{"Should succeed at upper bound", image.GravitySmart, nil},
		{"Should succeed in between bounds", image.GravityEast, nil},
		{"Should return error when below lower bound", image.GravityCentre - 1, ErrInvalidGravity},
		{"Should return error when above upper bound", image.GravitySmart + 1, ErrInvalidGravity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &GravityValidation{Gravity: tt.gravity}
			err := v.Validate()
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
