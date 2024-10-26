package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWidthValidation_ErrWidthSmall(t *testing.T) {
	w := WidthValidation{Width: 10, MinWidth: 20}
	err := w.Validate()
	assert.Equal(t, ErrWidthTooSmall, err)
}

func TestWidthValidation_ErrWidthLarge(t *testing.T) {
	w := WidthValidation{Width: 100, MaxWidth: 50}
	err := w.Validate()
	assert.Equal(t, ErrWidthTooLarge, err)
}

func TestWidthValidation_ValidWidth_Success(t *testing.T) {
	w := WidthValidation{Width: 30, MinWidth: 20, MaxWidth: 50}
	err := w.Validate()
	assert.Nil(t, err)
}
