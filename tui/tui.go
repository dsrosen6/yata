package tui

import (
	"fmt"

	"github.com/dsrosen6/yata/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	uncheckedIcon = "󰄱"
	checkedIcon   = "󰄵"
)

type ListHandler struct {
	*tview.List
}

func NewListHandler(initialTasks []*models.Task) *ListHandler {
	l := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedBackgroundColor(tcell.ColorDefault).
		SetSelectedTextColor(tcell.ColorBlue)
	l.SetBorder(true)
	l.SetTitle("tasks")
	l.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		}
		return event
	})

	for _, t := range initialTasks {
		l.AddItem(taskToListEntry(t), "", 0, nil)
	}

	return &ListHandler{
		List: l,
	}
}

func taskToListEntry(t *models.Task) string {
	return fmt.Sprintf("%s %s", checkbox(t), t.Title)
}

func checkbox(t *models.Task) string {
	if t.Complete {
		return checkedIcon
	}
	return uncheckedIcon
}
