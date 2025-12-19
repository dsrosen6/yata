package tui

import "github.com/dsrosen6/tea-flexbox/titlebox"

func (m *model) createTasksBox() titlebox.Box {
	boxStyle := allStyles.unfocusedBoxStyle
	titleStyle := allStyles.unfocusedBoxTitleStyle
	if m.currentFocus == focusTasks {
		boxStyle = allStyles.focusedBoxStyle
		titleStyle = allStyles.focusedBoxTitleStyle
	}

	return titlebox.New().
		SetTitle("[2]tasks").
		SetBody(m.taskList.View()).
		SetTitleAlignment(titlebox.AlignLeft).
		SetBoxStyle(boxStyle.Padding(0, 1)).
		SetTitleStyle(titleStyle)
}

func (m *model) createListsBox() titlebox.Box {
	boxStyle := allStyles.unfocusedBoxStyle
	titleStyle := allStyles.unfocusedBoxTitleStyle
	if m.currentFocus == focusLists {
		boxStyle = allStyles.focusedBoxStyle
		titleStyle = allStyles.focusedBoxTitleStyle
	}

	return titlebox.New().
		SetTitle("[1]lists").
		SetBody(m.listList.View()).
		SetTitleAlignment(titlebox.AlignLeft).
		SetBoxStyle(boxStyle.Padding(0, 1)).
		SetTitleStyle(titleStyle)
}

func (m *model) createTaskEntryBox() titlebox.Box {
	return titlebox.New().
		SetTitle("new task").
		SetBody(m.taskEntryForm.View()).
		SetTitleAlignment(titlebox.AlignLeft).
		SetBoxStyle(allStyles.focusedBoxStyle).
		SetTitleStyle(allStyles.focusedBoxTitleStyle)
}

func (m *model) createListEntryBox() titlebox.Box {
	return titlebox.New().
		SetTitle("new list").
		SetBody(m.listEntryForm.View()).
		SetTitleAlignment(titlebox.AlignLeft).
		SetBoxStyle(allStyles.focusedBoxStyle).
		SetTitleStyle(allStyles.focusedBoxTitleStyle)
}
