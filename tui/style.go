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
	baseFocused := lipgloss.NewStyle().Foreground(cfg.MainColor)
	baseUnfocused := lipgloss.NewStyle().Foreground(cfg.SecondaryColor)
	return styles{
		focusedBoxStyle:        baseFocused.Border(lipgloss.DoubleBorder()).BorderForeground(cfg.MainColor),
		focusedBoxTitleStyle:   baseFocused,
		focusedTaskStyle:       baseFocused,
		unfocusedBoxStyle:      baseUnfocused.Border(lipgloss.NormalBorder()).BorderForeground(cfg.SecondaryColor),
		unfocusedBoxTitleStyle: baseUnfocused,
		unfocusedTaskStyle:     baseUnfocused,
	}
}
