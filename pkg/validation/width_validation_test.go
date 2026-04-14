package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWidthValidation_ShouldReturnError_WhenWidthIsLessThanMinimum(t *testing.T) {
	w := WidthValidation{Width: 10, MinWidth: 20}
	err := w.Validate()
	assert.ErrorIs(t, err, ErrWidthTooSmall)
}

func TestWidthValidation_ShouldReturnError_WhenWidthExceedsMaximum(t *testing.T) {
	w := WidthValidation{Width: 100, MaxWidth: 50}
	err := w.Validate()
	assert.ErrorIs(t, err, ErrWidthTooLarge)
}

func TestWidthValidation_ShouldSucceed_WhenWidthIsWithinBounds(t *testing.T) {
	w := WidthValidation{Width: 30, MinWidth: 20, MaxWidth: 50}
	err := w.Validate()
	assert.NoError(t, err)
}
