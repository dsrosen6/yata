package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dsrosen6/yata/config"
)

type styles struct {
	focusedBoxStyle        lipgloss.Style
	focusedBoxTitleStyle   lipgloss.Style
	focusedTextStyle       lipgloss.Style
	unfocusedBoxStyle      lipgloss.Style
	unfocusedBoxTitleStyle lipgloss.Style
	unfocusedTextStyle     lipgloss.Style
	errorTextStyle         lipgloss.Style
}

func generateStyles(cfg *config.Config) styles {
	return styles{
		focusedBoxStyle:        lipgloss.NewStyle().Border(cfg.Focused.BorderType).BorderForeground(cfg.Focused.BorderColor).Foreground(cfg.Focused.TextColor),
		focusedBoxTitleStyle:   lipgloss.NewStyle().Foreground(cfg.Focused.BoxTitleColor),
		focusedTextStyle:       lipgloss.NewStyle().Foreground(cfg.Focused.TextColor),
		unfocusedBoxStyle:      lipgloss.NewStyle().Border(cfg.Unfocused.BorderType).BorderForeground(cfg.Unfocused.BorderColor),
		unfocusedBoxTitleStyle: lipgloss.NewStyle().Foreground(cfg.Unfocused.BoxTitleColor),
		unfocusedTextStyle:     lipgloss.NewStyle().Foreground(cfg.Unfocused.TextColor),
		errorTextStyle:         lipgloss.NewStyle().Foreground(cfg.ErrorTextColor),
	}
}
