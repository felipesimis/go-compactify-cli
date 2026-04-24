package ui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestWarn(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "should render warning with message",
			message: "This is a warning",
		},
		{
			name:    "should handle empty message",
			message: "",
		},
		{
			name:    "should handle message with leading and trailing spaces",
			message: "  This is a warning with spaces   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Warn(tt.message)
			cleanResult := lipgloss.Sprint(result)

			assert.Contains(t, cleanResult, "⚠️")
			assert.Contains(t, cleanResult, tt.message)
			assert.Contains(t, cleanResult, "┃")

			iconIndex := strings.Index(cleanResult, "⚠️")
			borderIndex := strings.Index(cleanResult, tt.message)

			if tt.message != "" {
				assert.True(t, iconIndex < borderIndex, "Icon should be before the message in the output")
			}
		})
	}
}
