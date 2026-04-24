package ui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var (
	warningBorderColor = lipgloss.Color("#F59E0B")
	warningTextColor   = lipgloss.Color("#FEF3C7")

	errorBorderColor = lipgloss.Color("#EF4444")
	errorTextColor   = lipgloss.Color("#FEE2E2")

	calloutStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder(), false, false, false, true).
			Padding(0, 2).Margin(1, 0)

	iconStyle    = lipgloss.NewStyle().Bold(true).MarginRight(2)
	contentStyle = lipgloss.NewStyle()
)

func renderCallout(icon string, message string, borderColor color.Color, textColor color.Color) string {
	style := calloutStyle.BorderForeground(borderColor)
	iconPart := iconStyle.Foreground(textColor).Render(icon)
	contentPart := contentStyle.Foreground(textColor).Render(message)

	return style.Render(lipgloss.JoinHorizontal(lipgloss.Left, iconPart, contentPart))
}

func Warn(message string) string {
	return renderCallout("⚠️", message, warningBorderColor, warningTextColor)
}

func Error(message string) string {
	return renderCallout("❌", message, errorBorderColor, errorTextColor)
}
