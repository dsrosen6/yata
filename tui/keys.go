package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	quit            key.Binding
	focusPanelLeft  key.Binding
	focusPanelRight key.Binding
	focusProjects   key.Binding
	focusTasks      key.Binding

	cancelEntry key.Binding
	delete      key.Binding // universal, depends on focused list

	newItem            key.Binding
	toggleTaskComplete key.Binding
}

var defaultKeyMap = keyMap{
	quit: key.NewBinding(
		// ctrl+c only there for the habit, q is all that will show in help
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	focusPanelLeft: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("h/left", "focus panel left"),
	),
	focusPanelRight: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("l/right", "focus panel right"),
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
	newItem: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "new"),
	),
	toggleTaskComplete: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("space/enter", "toggle complete"),
	),
}
