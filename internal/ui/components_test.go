package ui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/suite"
)

type CalloutTestSuite struct {
	suite.Suite
}

func (s *CalloutTestSuite) TestCallouts() {
	components := []struct {
		name   string
		render func(string) string
		icon   string
	}{
		{"Warn", Warn, "⚠️"},
		{"Error", Error, "❌"},
	}

	scenarios := []struct {
		name  string
		input string
	}{
		{"should render callout with message", "This is a message"},
		{"should handle empty message", ""},
		{"should handle message with leading and trailing spaces", "  This is a message with spaces   "},
	}

	for _, comp := range components {
		for _, sc := range scenarios {
			testName := comp.name + "_" + sc.name
			s.Run(testName, func() {
				result := comp.render(sc.input)
				cleanResult := lipgloss.Sprint(result)

				s.Contains(cleanResult, comp.icon)
				s.Contains(cleanResult, sc.input)
				s.Contains(cleanResult, "┃")

				iconIndex := strings.Index(cleanResult, comp.icon)
				msgIndex := strings.Index(cleanResult, sc.input)

				if sc.input != "" {
					s.True(iconIndex < msgIndex, "Icon should be before the message in the output")
				}
			})
		}
	}
}

func TestCalloutTestSuite(t *testing.T) {
	suite.Run(t, new(CalloutTestSuite))
}
