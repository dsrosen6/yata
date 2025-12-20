package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/dsrosen6/yata/models"
)

func initialTaskList(tasks []*models.Task) list.Model {
	items := tasksToItems(tasks)
	ls := list.New(items, taskItemDelegate{}, 10, 10)
	ls.SetShowStatusBar(false)
	ls.SetShowTitle(false)
	ls.SetShowHelp(false)
	ls.SetFilteringEnabled(false)
	return ls
}

func initialProjectList(projects []*models.Project) list.Model {
	items := projectsToItems(projects)
	ls := list.New(items, projectItemDelegate{}, 10, 10)
	ls.SetShowStatusBar(false)
	ls.SetShowTitle(false)
	ls.SetShowHelp(false)
	ls.SetFilteringEnabled(false)
	return ls
}
