package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatValidation_ShouldReturnError_WhenFormatIsEmpty(t *testing.T) {
	f := FormatValidation{Format: ""}
	err := f.Validate()
	assert.ErrorIs(t, err, ErrFormatRequired)
}

func TestFormatValidation_ShouldReturnError_WhenFormatIsNotSupported(t *testing.T) {
	f := FormatValidation{Format: "gif"}
	err := f.Validate()
	assert.ErrorIs(t, err, ErrInvalidFormat)
}

func TestFormatValidation_ShouldSucceed_WhenFormatIsSupported(t *testing.T) {
	f := FormatValidation{Format: "webp"}
	err := f.Validate()
	assert.NoError(t, err)
}
