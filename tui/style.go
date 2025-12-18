package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dsrosen6/yata/config"
)

type styles struct {
	focusedBoxStyle        lipgloss.Style
	focusedBoxTitleStyle   lipgloss.Style
	focusedTaskStyle       lipgloss.Style
	unfocusedBoxStyle      lipgloss.Style
	unfocusedBoxTitleStyle lipgloss.Style
	unfocusedTaskStyle     lipgloss.Style
}

func generateStyles(cfg *config.Config) styles {
	return styles{
		focusedBoxStyle:        lipgloss.NewStyle().Border(cfg.Focused.BorderType).BorderForeground(cfg.Focused.BorderColor).Foreground(cfg.Focused.TextColor),
		focusedBoxTitleStyle:   lipgloss.NewStyle().Foreground(cfg.Focused.BoxTitleColor),
		focusedTaskStyle:       lipgloss.NewStyle().Foreground(cfg.Focused.TextColor),
		unfocusedBoxStyle:      lipgloss.NewStyle().Border(cfg.Unfocused.BorderType).BorderForeground(cfg.Unfocused.BorderColor),
		unfocusedBoxTitleStyle: lipgloss.NewStyle().Foreground(cfg.Unfocused.BoxTitleColor),
		unfocusedTaskStyle:     lipgloss.NewStyle().Foreground(cfg.Unfocused.TextColor),
	}
}
