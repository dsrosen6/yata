package tui

import (
	"github.com/charmbracelet/bubbles/list"
)

func initialTaskList() list.Model {
	ls := list.New([]list.Item{}, taskItemDelegate{}, 10, 10)
	ls.SetShowStatusBar(false)
	ls.SetShowTitle(false)
	ls.SetShowHelp(false)
	ls.SetFilteringEnabled(false)
	return ls
}

func initialProjectList() list.Model {
	ls := list.New([]list.Item{}, projectItemDelegate{maxWidth: 20}, 10, 10)
	ls.SetShowStatusBar(false)
	ls.SetShowTitle(false)
	ls.SetShowHelp(false)
	ls.SetFilteringEnabled(false)
	return ls
}
