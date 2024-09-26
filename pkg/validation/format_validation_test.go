package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatValidation_Validate_EmptyFormat(t *testing.T) {
	f := FormatValidation{Format: ""}
	err := f.Validate()
	assert.Equal(t, ErrFormatRequired, err)
}

func TestFormatValidation_Validate_InvalidFormat(t *testing.T) {
	f := FormatValidation{Format: "invalid"}
	err := f.Validate()
	assert.Equal(t, ErrInvalidFormat, err)
}

func TestFormatValidation_Validate_Success(t *testing.T) {
	f := FormatValidation{Format: "jpg"}
	err := f.Validate()
	assert.Nil(t, err)
}
