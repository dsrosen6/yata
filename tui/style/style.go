package style

import "github.com/charmbracelet/lipgloss"

var DefaultColor = lipgloss.Color("205")

func BorderStyle(color string) lipgloss.Style {
	c := DefaultColor
	if color != "" {
		c = lipgloss.Color(color)
	}

	return lipgloss.NewStyle().BorderForeground(c)
}

func FocusedStyle(color string) lipgloss.Style {
	c := DefaultColor
	if color != "" {
		c = lipgloss.Color(color)
	}

	return lipgloss.NewStyle().Foreground(c)
}

func UnfocusedStyle(color string) lipgloss.Style {
	var c lipgloss.Color
	if color != "" {
		c = lipgloss.Color(color)
	}

	s := lipgloss.NewStyle().Foreground(c)
	return s
}
