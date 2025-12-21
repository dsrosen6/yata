package tui

type focus int

const (
	focusTasks focus = iota
	focusProjects
	focusTaskEntry
	focusProjectEntry
)

func (f focus) isEntry() bool {
	return f == focusTaskEntry || f == focusProjectEntry
}
