package validation

import (
	"testing"

	"github.com/h2non/bimg"
	"github.com/stretchr/testify/assert"
)

func TestGravityValidation_ShouldReturnError_WhenGravityIsBelowCentre(t *testing.T) {
	v := &GravityValidation{Gravity: bimg.GravityCentre - 1}
	err := v.Validate()
	assert.ErrorIs(t, err, ErrInvalidGravity)
}

func TestGravityValidation_ShouldReturnError_WhenGravityIsAboveSmart(t *testing.T) {
	v := &GravityValidation{Gravity: bimg.GravitySmart + 1}
	err := v.Validate()
	assert.ErrorIs(t, err, ErrInvalidGravity)
}

func TestGravityValidation_ShouldSucceed_WhenGravityIsAtCentre(t *testing.T) {
	v := &GravityValidation{Gravity: bimg.GravityCentre}
	err := v.Validate()
	assert.NoError(t, err)
}

func TestGravityValidation_ShouldSucceed_WhenGravityIsAtSmart(t *testing.T) {
	v := &GravityValidation{Gravity: bimg.GravitySmart}
	err := v.Validate()
	assert.NoError(t, err)
}

func TestGravityValidation_ShouldSucceed_WhenGravityIsSomewhereInBetween(t *testing.T) {
	v := &GravityValidation{Gravity: bimg.GravityCentre + 1}
	err := v.Validate()
	assert.NoError(t, err)
}
