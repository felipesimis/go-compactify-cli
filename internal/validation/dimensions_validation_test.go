package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDimensionsValidation_ShouldReturnError_WhenWidthIsLessThanMin(t *testing.T) {
	v := &DimensionsValidation{Width: 0, Height: 100}
	err := v.Validate()
	assert.ErrorIs(t, err, ErrInvalidDimensions)
}

func TestDimensionsValidation_ShouldReturnError_WhenHeightIsLessThanMin(t *testing.T) {
	v := &DimensionsValidation{Width: 100, Height: 0}
	err := v.Validate()
	assert.ErrorIs(t, err, ErrInvalidDimensions)
}

func TestDimensionsValidation_ShouldReturnError_WhenBothDimensionsAreInvalid(t *testing.T) {
	v := &DimensionsValidation{Width: -10, Height: -10}
	err := v.Validate()
	assert.ErrorIs(t, err, ErrInvalidDimensions)
}

func TestDimensionsValidation_ShouldSucceed_WhenDimensionsAreValid(t *testing.T) {
	v := &DimensionsValidation{Width: 10, Height: 10}
	err := v.Validate()
	assert.NoError(t, err)
}
