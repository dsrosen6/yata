package tui

import "github.com/dsrosen6/yata/tui/render/titlebox"

func (m *model) createProjectsBox() titlebox.Box {
	boxStyle := allStyles.unfocusedBoxStyle
	titleStyle := allStyles.unfocusedBoxTitleStyle
	if m.currentFocus == focusProjects {
		boxStyle = allStyles.focusedBoxStyle
		titleStyle = allStyles.focusedBoxTitleStyle
	}

	border := boxStyle.GetBorderStyle().Top
	title := "[1]" + border + "projects"

	return titlebox.New().
		SetTitle(title).
		SetBody(m.projectList.View()).
		SetTitleAlignment(titlebox.AlignLeft).
		SetBoxStyle(boxStyle.Padding(0, 1)).
		SetTitleStyle(titleStyle)
}

func (m *model) createTasksBox() titlebox.Box {
	boxStyle := allStyles.unfocusedBoxStyle
	titleStyle := allStyles.unfocusedBoxTitleStyle
	if m.currentFocus == focusTasks {
		boxStyle = allStyles.focusedBoxStyle
		titleStyle = allStyles.focusedBoxTitleStyle
	}

	border := boxStyle.GetBorderStyle().Top
	title := "[2]" + border + "tasks"

	return titlebox.New().
		SetTitle(title).
		SetBody(m.taskList.View()).
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

func (m *model) createProjectEntryBox() titlebox.Box {
	return titlebox.New().
		SetTitle("new project").
		SetBody(m.projectEntryForm.View()).
		SetTitleAlignment(titlebox.AlignLeft).
		SetBoxStyle(allStyles.focusedBoxStyle).
		SetTitleStyle(allStyles.focusedBoxTitleStyle)
}
