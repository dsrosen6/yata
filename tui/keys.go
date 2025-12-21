package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	quit          key.Binding
	focusProjects key.Binding
	focusTasks    key.Binding

	cancelEntry key.Binding
	delete      key.Binding // universal, depends on focused list

	newProject         key.Binding
	newTask            key.Binding
	toggleTaskComplete key.Binding
}

var defaultKeyMap = keyMap{
	quit: key.NewBinding(
		// ctrl+c only there for the habit, q is all that will show in help
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),

	focusProjects: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "focus projects"),
	),
	focusTasks: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "focus tasks"),
	),
	cancelEntry: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	delete: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "delete"),
	),
	newProject: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "new project"),
	),
	newTask: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "new task"),
	),
	toggleTaskComplete: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("space/enter", "toggle complete"),
	),
}
