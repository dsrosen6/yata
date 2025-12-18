package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dsrosen6/tea-flexbox/titlebox"
	"github.com/dsrosen6/yata/config"
)

func mainStyle(cfg *config.Config) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(cfg.MainColor).
		Foreground(cfg.MainColor).
		Align(lipgloss.Center, lipgloss.Center)
}

func boxTitleStyle(cfg *config.Config) lipgloss.Style {
	return titlebox.DefaultTitleStyle.Foreground(cfg.MainColor)
}
