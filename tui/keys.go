package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	navigationKeys
	entryKeys
}

type navigationKeys struct {
	quit               key.Binding
	toggleHelp         key.Binding
	focusProjects      key.Binding
	focusTasks         key.Binding
	delete             key.Binding
	newTask            key.Binding
	newProject         key.Binding
	toggleTaskComplete key.Binding
}

type entryKeys struct {
	cancelEntry key.Binding
	submit      key.Binding
}

var (
	defaultKeyMap = keyMap{
		navigationKeys: defaultNavKeys,
		entryKeys:      defaultEntryKeys,
	}
	helpStyle = lipgloss.NewStyle().Padding(0, 1).AlignHorizontal(lipgloss.Center)
)

var defaultNavKeys = navigationKeys{
	quit: key.NewBinding(
		// ctrl+c only there for the habit, q is all that will show in help
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	toggleHelp: key.NewBinding(
		key.WithKeys("H"),
	),
	focusProjects: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "focus projects"),
	),
	focusTasks: key.NewBinding(
		key.WithKeys("2", "focus tasks"),
		key.WithHelp("2", "focus tasks"),
	),
	delete: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "delete"),
	),
	newTask: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new task"),
	),
	newProject: key.NewBinding(
		key.WithKeys("N"),
		key.WithHelp("N", "new project"),
	),
	toggleTaskComplete: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "complete"),
	),
}

var defaultEntryKeys = entryKeys{
	cancelEntry: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
}

func (k navigationKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.focusProjects, k.focusTasks}
}

func (k navigationKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.focusProjects, k.focusTasks},
	}
}

func (m *model) helpKeys() []key.Binding {
	if m.currentFocus.isEntry() {
		return []key.Binding{m.keys.cancelEntry, m.keys.submit}
	}

	k := []key.Binding{m.keys.newTask, m.keys.newProject}
	switch m.currentFocus {
	case focusProjects:
		if m.selectedProjectID() != nil {
			k = append(k, m.keys.delete)
		}
	case focusTasks:
		if m.selectedTaskID() != 0 {
			tc := taskCompleteHelp(m.keys.toggleTaskComplete, m.selectedTask().Complete)
			k = append(k, tc, m.keys.delete)
		}
	}

	return k
}

// taskCompleteHelp takes an original key binding for task completion, and creates
// a new one with dynamic help text. This is for use in help outputs; the original
// should still be used for actual key detection in Update.
func taskCompleteHelp(original key.Binding, complete bool) key.Binding {
	c := "complete"
	if complete {
		c = "uncomplete"
	}

	return key.NewBinding(
		key.WithKeys(original.Keys()...),
		key.WithHelp(original.Help().Key, c),
	)
}
