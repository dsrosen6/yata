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
		focusedBoxStyle:        lipgloss.NewStyle().Border(cfg.Style.Focused.BorderType).BorderForeground(cfg.Style.Focused.BorderColor).Foreground(cfg.Style.Focused.TextColor),
		focusedBoxTitleStyle:   lipgloss.NewStyle().Foreground(cfg.Style.Focused.BoxTitleColor),
		focusedTextStyle:       lipgloss.NewStyle().Foreground(cfg.Style.Focused.TextColor),
		unfocusedBoxStyle:      lipgloss.NewStyle().Border(cfg.Style.Unfocused.BorderType).BorderForeground(cfg.Style.Unfocused.BorderColor),
		unfocusedBoxTitleStyle: lipgloss.NewStyle().Foreground(cfg.Style.Unfocused.BoxTitleColor),
		unfocusedTextStyle:     lipgloss.NewStyle().Foreground(cfg.Style.Unfocused.TextColor),
		errorTextStyle:         lipgloss.NewStyle().Foreground(cfg.Style.ErrorTextColor),
	}
}
