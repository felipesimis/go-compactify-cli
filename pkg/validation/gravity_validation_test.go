package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGravityValidation_Validate_Error(t *testing.T) {
	v := &GravityValidation{Gravity: -5}
	err := v.Validate()
	v2 := &GravityValidation{Gravity: 10}
	err2 := v2.Validate()
	assert.Equal(t, ErrInvalidGravity, err)
	assert.Equal(t, ErrInvalidGravity, err2)
}

func TestGravityValidation_Validate_Success(t *testing.T) {
	v := &GravityValidation{Gravity: 0}
	err := v.Validate()
	v2 := &GravityValidation{Gravity: 5}
	err2 := v2.Validate()
	assert.Nil(t, err)
	assert.Nil(t, err2)
}
