package tui

import tea "github.com/charmbracelet/bubbletea"

type (
	focus          int
	changeFocusMsg struct{ focus focus }
)

const (
	focusTasks focus = iota
	focusProjects
	focusTaskEntry
	focusProjectEntry
)

func (f focus) isEntry() bool {
	return f == focusTaskEntry || f == focusProjectEntry
}

func changeFocus(f focus) tea.Cmd {
	return func() tea.Msg {
		return changeFocusMsg{focus: f}
	}
}
