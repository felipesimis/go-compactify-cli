package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGravity_IsValid_ShouldEnforceDomainBounds(t *testing.T) {
	tests := []struct {
		name     string
		input    Gravity
		expected bool
	}{
		{"Valid lower bound", GravityCentre, true},
		{"Valid upper bound", GravitySmart, true},
		{"Valid middle value", GravityEast, true},
		{"Invalid negative value", Gravity(-1), false},
		{"Invalid above max bound", maxGravity, false},
		{"Invalid arbitrary high value", Gravity(99), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}
