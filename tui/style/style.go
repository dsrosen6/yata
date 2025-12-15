package style

import "github.com/charmbracelet/lipgloss"

var (
	DefaultColor        lipgloss.Color // default terminal color
	DefaultFocusedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(DefaultColor))
)

func BorderStyle(color string) lipgloss.Style {
	c := DefaultColor
	if color != "" {
		c = lipgloss.Color(color)
	}

	return lipgloss.NewStyle().BorderForeground(c)
}

func FocusedStyle(color string) lipgloss.Style {
	s := DefaultFocusedStyle
	if color != "" {
		s = s.Foreground(lipgloss.Color(color))
	}

	return s
}

func UnfocusedStyle(focusStyle lipgloss.Style, overrideColor string) lipgloss.Style {
	if overrideColor != "" {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(overrideColor))
	}

	return focusStyle.Faint(true)
}
