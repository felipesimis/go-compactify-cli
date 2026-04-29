package ui

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

func renderCallout(icon string, message string, borderColor color.Color, textColor color.Color) string {
	style := styleCalloutBase.BorderForeground(borderColor)
	iconPart := styleCalloutIcon.Foreground(textColor).Render(icon)
	contentPart := styleCalloutText.Foreground(textColor).Render(message)

	return style.Render(lipgloss.JoinHorizontal(lipgloss.Left, iconPart, contentPart))
}

func Warn(message string) string {
	return renderCallout("⚠️", message, colorWarnBorder, colorWarnText)
}

func Error(message string) string {
	return renderCallout("❌", message, colorErrorBorder, colorErrText)
}

func Success(message string) string {
	return renderCallout("✅", message, colorSuccessBorder, colorSuccessText)
}

type Item struct {
	Label         string
	Value         string
	IsHighlighted bool
}

type Panel struct {
	Title string
	Items []Item
}

func RenderPanel(p Panel) string {
	var lines []string

	lines = append(lines, styleTitle.Render(p.Title))
	for _, item := range p.Items {
		valRendered := styleValue.Render(item.Value)
		if item.IsHighlighted {
			valRendered = styleHero.Render(item.Value)
		}
		lines = append(lines, styleLabel.Render(item.Label)+valRendered)
	}
	return lipgloss.NewStyle().Width(30).Render(strings.Join(lines, "\n"))
}

func RenderDashboard(left Panel, right Panel, footerTitle, footerLine string) string {
	body := lipgloss.JoinHorizontal(lipgloss.Top, RenderPanel(left), RenderPanel(right))
	if footerTitle == "" {
		return styleBox.Render(body)
	}

	width := lipgloss.Width(body)
	content := strings.Join([]string{
		body,
		"",
		renderFooter(footerTitle, footerLine, width),
	}, "\n")

	return styleBox.Render(content)
}

func renderFooter(footerTitle, footerLine string, width int) string {
	divider := styleDivider.Render(strings.Repeat("─", width))
	title := styleFooterTitle.Render(footerTitle)
	line := styleFooterLine.Render(footerLine)
	return strings.Join([]string{divider, title, line}, "\n")
}

func RenderErrorList(errs []error) string {
	errCount := len(errs)
	if errCount == 0 {
		return ""
	}

	var sb strings.Builder
	headerText := fmt.Sprintf(" %d ERRORS DETECTED ", errCount)
	sb.WriteString(styleErrorHeader.Render(headerText) + "\n")

	for _, err := range errs {
		msg := err.Error()
		parts := strings.SplitN(msg, ": ", 2)

		var formattedMsg string
		if len(parts) == 2 {
			formattedMsg = styleErrorPath.Render(parts[0]) + ": " + styleErrorMsg.Render(parts[1])
		} else {
			formattedMsg = styleErrorMsg.Render(msg)
		}

		line := lipgloss.JoinHorizontal(lipgloss.Left, styleErrorSymbol.Render("❌"), formattedMsg)
		sb.WriteString(styleErrorItem.Render(line) + "\n")
	}
	return sb.String()
}
