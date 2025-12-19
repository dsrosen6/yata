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

func initialListList(lists []*models.List) list.Model {
	items := listsToItems(lists)
	ls := list.New(items, listItemDelegate{}, 10, 10)
	ls.SetShowStatusBar(false)
	ls.SetShowTitle(false)
	ls.SetShowHelp(false)
	ls.SetFilteringEnabled(false)
	return ls
}
