package ui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var (
	warningColor = lipgloss.Color("#F59E0B")
	textColor    = lipgloss.Color("#FEF3C7")

	calloutStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder(), false, false, false, true).
			Padding(0, 2).Margin(1, 0)

	iconStyle    = lipgloss.NewStyle().Bold(true).MarginRight(2)
	contentStyle = lipgloss.NewStyle().Foreground(textColor)
)

func renderCallout(icon string, message string, color color.Color) string {
	style := calloutStyle.BorderForeground(color)
	iconPart := iconStyle.Foreground(color).Render(icon)
	contentPart := contentStyle.Render(message)

	return style.Render(lipgloss.JoinHorizontal(lipgloss.Left, iconPart, contentPart))
}

func Warn(message string) string {
	return renderCallout("⚠️", message, warningColor)
}
