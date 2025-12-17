package tui

import (
	"context"
	"fmt"

	"github.com/dsrosen6/yata/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	uncheckedIcon = "󰄱"
	checkedIcon   = "󰄵"
)

func (a *app) newTaskList() *tview.List {
	l := tview.NewList().
		ShowSecondaryText(false)
	setDefaultListColors(l)
	l.SetBorder(true)
	l.SetTitle("tasks")
	l.SetTitleAlign(tview.AlignLeft)
	l.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case 'd':
			if len(a.tasks) == 0 {
				return event
			}
			sel := a.tasks[l.GetCurrentItem()]
			if err := a.deleteTask(context.Background(), sel.ID); err != nil {
				// TODO: handle this
				return event
			}
		case ' ':
			// check or uncheck
			sel := a.tasks[l.GetCurrentItem()]
			t := *sel
			t.Complete = !t.Complete
			if err := a.updateTask(context.Background(), &t); err != nil {
				// TODO: handle this
				return event
			}
		}
		return event
	})

	return l
}

func (a *app) initTaskList(lh *tview.List) error {
	if err := a.refreshListTasks(lh); err != nil {
		return fmt.Errorf("refreshing tasks: %w", err)
	}

	return nil
}

func (a *app) refreshListTasks(lh *tview.List) error {
	sel := lh.GetCurrentItem()
	lh.Clear()
	for _, t := range a.tasks {
		lh.AddItem(taskToListEntry(t), "", 0, nil)
	}
	if len(a.tasks) != 0 {
		lh.SetCurrentItem(sel)
	}

	return nil
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
