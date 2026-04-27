package ui

import (
	"errors"
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
				rawResult := comp.render(sc.input)
				cleanResult := lipgloss.Sprint(rawResult)

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

type PanelTestSuite struct {
	suite.Suite
}

func (s *PanelTestSuite) TestRenderPanel_ShouldIncludeAllText_WhenProvidedWithValidPanel() {
	panel := Panel{
		Title: "Test Title",
		Items: []Item{{"Label", "Value", false}},
	}

	cleanResult := lipgloss.Sprint(RenderPanel(panel))

	s.Contains(cleanResult, "Test Title")
	s.Contains(cleanResult, "Label")
	s.Contains(cleanResult, "Value")
}

func (s *PanelTestSuite) TestRenderPanel_ShouldHighlightValue_WhenItemIsHighlighted() {
	panel := Panel{
		Title: "Test Title",
		Items: []Item{{"Label", "Value", true}},
	}

	rawResult := RenderPanel(panel)
	s.Contains(rawResult, styleHero.Render("Value"), "Highlighted value should be rendered with hero style")
}

func (s *PanelTestSuite) TestRenderPanel_ShouldApplyValueStyle_WhenItemIsNotHighlighted() {
	panel := Panel{
		Items: []Item{{"Label", "Value", false}},
	}

	rawResult := RenderPanel(panel)
	s.Contains(rawResult, styleValue.Render("Value"), "Non-highlighted value should be rendered with value style")
}

func TestPanelTestSuite(t *testing.T) {
	suite.Run(t, new(PanelTestSuite))
}

type DashboardTestSuite struct {
	suite.Suite
}

func (s *DashboardTestSuite) TestRenderDashboard_ShouldAlignPanelsHorizontally_WhenTwoPanelsAreProvided() {
	left := Panel{Title: "Left Panel", Items: []Item{{"L1", "Value1", false}}}
	right := Panel{Title: "Right Panel", Items: []Item{{"R1", "Value2", false}}}

	rawResult := RenderDashboard(left, right, "", "")
	cleanResult := lipgloss.Sprint(rawResult)

	s.Contains(cleanResult, "Left Panel")
	s.Contains(cleanResult, "Right Panel")

	lines := strings.Split(cleanResult, "\n")
	foundLeft := false
	foundRightInSameLine := false

	for _, line := range lines {
		if strings.Contains(line, "Left Panel") {
			foundLeft = true
			if strings.Contains(line, "Right Panel") {
				foundRightInSameLine = true
			}
		}
	}

	s.True(foundLeft, "Left panel title should be present in the output")
	s.True(foundRightInSameLine, "Right panel title should be on the same line as the left panel")
}

func (s *DashboardTestSuite) TestRenderDashboard_ShouldApplyBoxStyle_WhenRendered() {
	left := Panel{Title: "L"}
	right := Panel{Title: "R"}

	rawResult := RenderDashboard(left, right, "", "")

	s.Contains(rawResult, "╭", "Output should contain box style top-left corner")
	s.Contains(rawResult, "╯", "Output should contain box style bottom-right corner")
}

func (s *DashboardTestSuite) TestRenderDashboard_ShouldRenderPanelsInBox_WhenNoFooterProvided() {
	left := Panel{Title: "L"}
	right := Panel{Title: "R"}

	cleanResult := lipgloss.Sprint(RenderDashboard(left, right, "", "Footer line"))

	s.Contains(cleanResult, "L", "Left panel content should be present in the output")
	s.Contains(cleanResult, "R", "Right panel content should be present in the output")

	s.NotContains(cleanResult, "Footer line", "Footer line should not be present when footer title is empty")
}

func (s *DashboardTestSuite) TestRenderDashboard_ShouldRenderFooter_WhenFooterTitleProvided() {
	left := Panel{Title: "L"}
	right := Panel{Title: "R"}
	footerTitle := "Footer Title"
	footerLine := "Footer Line"

	cleanResult := lipgloss.Sprint(RenderDashboard(left, right, footerTitle, footerLine))

	s.Contains(cleanResult, footerTitle, "Footer title should be present in the output")
	s.Contains(cleanResult, footerLine, "Footer line should be present in the output")
}

func (s *DashboardTestSuite) TestRenderDashboard_ShouldApplyFooterStyles_WhenFooterIsProvided() {
	left := Panel{Title: "L"}
	right := Panel{Title: "R"}
	footerTitle := "Footer Title"
	footerLine := "Footer Line"

	rawResult := RenderDashboard(left, right, footerTitle, footerLine)
	cleanResult := lipgloss.Sprint(rawResult)

	s.Contains(cleanResult, "─", "Footer should contain a horizontal divider line")

	expectedTitle := styleFooterTitle.Render(footerTitle)
	s.Contains(rawResult, strings.TrimSpace(expectedTitle), "Footer title should have correct ANSI styles")

	s.Contains(cleanResult, footerTitle, "Footer title should be present in the output")
	s.Contains(cleanResult, footerLine, "Footer line should be present in the output")
}

func TestDashboardTestSuite(t *testing.T) {
	suite.Run(t, new(DashboardTestSuite))
}

type ErrorListTestSuite struct {
	suite.Suite
}

func (s *ErrorListTestSuite) TestRenderErrorList_ShouldReturnEmpty_WhenNoErrors() {
	s.Empty(RenderErrorList([]error{}))
}

func (s *ErrorListTestSuite) TestRenderErrorList_ShouldIncludeStyledHeader_WhenErrorsExist() {
	errors := []error{errors.New("Error 1")}

	rawResult := RenderErrorList(errors)
	cleanResult := lipgloss.Sprint(rawResult)

	expectedHeader := styleErrorHeader.Render(" 1 ERRORS DETECTED ")
	s.Contains(cleanResult, "1 ERRORS DETECTED", "Header should indicate the number of errors detected")
	s.Contains(rawResult, expectedHeader, "Header should have correct styles applied")
}

func (s *ErrorListTestSuite) TestRenderErrorList_ShouldIncludeErrorItems_WhenMultipleErrorsExist() {
	errors := []error{
		errors.New("error 1"),
		errors.New("error 2"),
	}

	rawResult := RenderErrorList(errors)
	cleanResult := lipgloss.Sprint(rawResult)

	s.Contains(cleanResult, "❌ error 1")
	s.Contains(cleanResult, "❌ error 2")
	expectedSymbol := styleErrorSymbol.Render("❌")
	s.Contains(rawResult, expectedSymbol, "Each error item should have the error symbol with correct styles")
}

func TestErrorListTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorListTestSuite))
}
