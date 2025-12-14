package input

import "github.com/charmbracelet/lipgloss"

var defaultFocusedColor = lipgloss.Color("205")

func focusedStyle(color string) lipgloss.Style {
	c := defaultFocusedColor
	if color != "" {
		c = lipgloss.Color(color)
	}

	return lipgloss.NewStyle().Foreground(c)
}

func unfocusedStyle(color string) lipgloss.Style {
	s := lipgloss.NewStyle()
	if color != "" {
		s = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	}

	return s
}
