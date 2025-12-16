package tui

import (
	"context"

	"github.com/dsrosen6/yata/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newTaskEntryBox() *tview.InputField {
	f := tview.NewInputField().
		SetLabel("Title: ").
		SetFieldBackgroundColor(tcell.ColorDefault)
	f.SetBorder(true)

	f.SetDoneFunc(func(key tcell.Key) {
		ctx := context.Background()
		title := f.GetText()
		task := &models.Task{
			Title: title,
		}

		if _, err := app.addTask(ctx, task); err != nil {
			return // TODO: do something
		}

		if err := app.refreshTasks(ctx); err != nil {
			return // TODO: do something
		}

		app.AddingTask = false
		app.SetFocus(app.ListFlex)
		app.ListFlex.RemoveItem(f)
	})

	return f
}
