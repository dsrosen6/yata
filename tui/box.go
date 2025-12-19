package tui

import "github.com/dsrosen6/tea-flexbox/titlebox"

func (m *model) createTasksBox() titlebox.Box {
	boxStyle := m.unfocusedBoxStyle
	titleStyle := m.unfocusedBoxTitleStyle
	if m.currentFocus == focusTasks {
		boxStyle = m.focusedBoxStyle
		titleStyle = m.focusedBoxTitleStyle
	}

	return titlebox.New().
		SetTitle("tasks").
		SetBody(m.taskList.View()).
		SetTitleAlignment(titlebox.AlignLeft).
		SetBoxStyle(boxStyle.Padding(0, 1)).
		SetTitleStyle(titleStyle)
}

func (m *model) createListsBox() titlebox.Box {
	boxStyle := m.unfocusedBoxStyle
	titleStyle := m.unfocusedBoxTitleStyle
	if m.currentFocus == focusLists {
		boxStyle = m.focusedBoxStyle
		titleStyle = m.focusedBoxTitleStyle
	}

	return titlebox.New().
		SetTitle("lists").
		SetBody(m.listList.View()).
		SetTitleAlignment(titlebox.AlignLeft).
		SetBoxStyle(boxStyle.Padding(0, 1)).
		SetTitleStyle(titleStyle)
}

func (m *model) createTaskEntryBox() titlebox.Box {
	return titlebox.New().
		SetTitle("new task").
		SetBody(m.entryForm.View()).
		SetTitleAlignment(titlebox.AlignLeft).
		SetBoxStyle(m.focusedBoxStyle).
		SetTitleStyle(m.focusedBoxTitleStyle)
}

func (m *model) createListEntryBox() titlebox.Box {
	return titlebox.New().
		SetTitle("new list").
		SetTitleAlignment(titlebox.AlignLeft).
		SetBoxStyle(m.focusedBoxStyle).
		SetTitleStyle(m.focusedBoxTitleStyle)
}
