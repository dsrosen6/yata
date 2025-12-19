package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/models"
)

type (
	taskItem         struct{ *models.Task }
	taskListItem     struct{ *models.List }
	taskItemDelegate struct{}
	listItemDelegate struct{}
)

func (d taskItemDelegate) Height() int {
	return 1
}

func (d taskItemDelegate) Spacing() int {
	return 0
}

func (d taskItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d taskItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(taskItem)
	if !ok {
		return
	}

	checked := "󰄱"
	if i.Complete {
		checked = "󰄵"
	}

	str := fmt.Sprintf("%s %s", checked, i.Title)
	fn := allStyles.unfocusedTextStyle.Render
	if index == m.Index() {
		fn = allStyles.focusedTextStyle.Render
	}

	_, _ = fmt.Fprint(w, fn(str))
}

func (t taskItem) FilterValue() string {
	return t.Title
}

func (d listItemDelegate) Height() int {
	return 1
}

func (d listItemDelegate) Spacing() int {
	return 0
}

func (d listItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d listItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(taskListItem)
	if !ok {
		return
	}

	prepend := " "
	fn := allStyles.unfocusedTextStyle.Render
	if index == m.Index() {
		prepend = ">"
		fn = allStyles.focusedTextStyle.Render
	}
	str := fmt.Sprintf("%s%s", prepend, i.Title)
	_, _ = fmt.Fprint(w, fn(str))
}

func (l taskListItem) FilterValue() string {
	return l.Title
}
