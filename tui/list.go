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

type ListHandler struct {
	*tview.List
}

func NewListHandler() *ListHandler {
	l := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedBackgroundColor(tcell.ColorDefault).
		SetSelectedTextColor(tcell.ColorBlue)
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
			if len(app.Tasks) == 0 {
				return event
			}
			sel := app.Tasks[l.GetCurrentItem()]
			if err := app.deleteTask(context.Background(), sel.ID); err != nil {
				// TODO: handle this
				return event
			}
		case ' ':
			// check or uncheck
			sel := app.Tasks[l.GetCurrentItem()]
			t := *sel
			t.Complete = !t.Complete
			if err := app.updateTask(context.Background(), &t); err != nil {
				// TODO: handle this
				return event
			}
		}
		return event
	})

	return &ListHandler{
		List: l,
	}
}

func (lh *ListHandler) Init(ctx context.Context) error {
	if err := lh.RefreshTasks(ctx); err != nil {
		return fmt.Errorf("refreshing tasks: %w", err)
	}

	return nil
}

func (lh *ListHandler) RefreshTasks(ctx context.Context) error {
	sel := lh.GetCurrentItem()
	lh.Clear()
	for _, t := range app.Tasks {
		lh.AddItem(taskToListEntry(t), "", 0, nil)
	}
	if len(app.Tasks) != 0 {
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
