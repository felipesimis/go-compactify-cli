package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDimensionsValidation_Validate_Error(t *testing.T) {
	v := &DimensionsValidation{Width: -100, Height: -100}
	err := v.Validate()
	assert.Equal(t, ErrInvalidDimensions, err)
}
