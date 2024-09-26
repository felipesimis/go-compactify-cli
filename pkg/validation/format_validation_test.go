package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatValidation_Validate_EmptyFormat(t *testing.T) {
	f := FormatValidation{Format: ""}
	err := f.Validate()
	assert.Equal(t, ErrInvalidFormat, err)
}
