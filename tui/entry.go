package tui

import (
	"context"

	"github.com/dsrosen6/yata/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *app) newTaskEntryField() *tview.InputField {
	t := tview.NewInputField().SetLabel("Title: ")
	setDefaultInputColors(t)
	t.SetBorder(true)
	t.SetDoneFunc(a.handleTaskEntryDone)

	return t
}

func (a *app) handleTaskEntryDone(key tcell.Key) {
	if key == tcell.KeyEnter {
		ctx := context.Background()
		title := a.taskEntryField.GetText()
		task := &models.Task{
			Title: title,
		}

		if err := a.addTask(ctx, task); err != nil {
			return // TODO: do something
		}
	}
	a.addingTask = false
	a.SetFocus(a.mainFlex)
	a.rootFlex.RemoveItem(a.entryFlex)
	a.taskEntryField.SetText("")
}
