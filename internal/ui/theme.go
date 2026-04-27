package ui

import "charm.land/lipgloss/v2"

var (
	colorAccent    = lipgloss.Color("#10B981")
	colorSubtle    = lipgloss.Color("#3F3F46")
	colorDimText   = lipgloss.Color("#A1A1AA")
	colorWhiteText = lipgloss.Color("#F4F4F5")

	colorErrBg = lipgloss.Color("#7F1D1D")

	colorWarnBorder = lipgloss.Color("#F59E0B")
	colorWarnText   = lipgloss.Color("#FEF3C7")

	colorBorderError = lipgloss.Color("#EF4444")
	colorErrText     = lipgloss.Color("#FEE2E2")

	styleCalloutBase = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder(), false, false, false, true).
				Padding(0, 2).
				Margin(1, 0)

	styleCalloutIcon = lipgloss.NewStyle().Bold(true).MarginRight(2)
	styleCalloutText = lipgloss.NewStyle()

	styleTitle = lipgloss.NewStyle().Foreground(colorAccent).Bold(true).MarginBottom(1)
	styleLabel = lipgloss.NewStyle().Foreground(colorDimText).Width(12)
	styleValue = lipgloss.NewStyle().Foreground(colorWhiteText).Bold(true)
	styleHero  = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)

	styleDivider     = lipgloss.NewStyle().Foreground(colorSubtle)
	styleFooterTitle = styleTitle.MarginTop(1)
	styleFooterLine  = lipgloss.NewStyle().Foreground(colorDimText)

	styleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSubtle).
			Padding(1, 2)

	styleErrorHeader = lipgloss.NewStyle().Background(colorErrBg).Foreground(colorWhiteText).Bold(true).Padding(0, 1).Margin(1, 0, 0, 2)
	styleErrorSymbol = lipgloss.NewStyle().Foreground(colorErrText).Margin(0, 1, 0, 4)
	styleErrorPath   = lipgloss.NewStyle().Foreground(colorDimText)
	styleErrorMsg    = lipgloss.NewStyle().Foreground(colorWhiteText).Bold(true)

	styleErrorItem = lipgloss.NewStyle().MarginLeft(1)
)
