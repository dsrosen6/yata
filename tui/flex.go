package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newListFlex(lh *ListHandler) *tview.Flex {
	f := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(lh, 0, 2, true)

	f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			if !app.AddingTask {
				app.AddingTask = true
				i := newTaskEntryBox()
				f.AddItem(i, 0, 1, true)
				app.SetFocus(i)
				return nil
			}
		}
		return event
	})
	return f
}

func newSummaryFlex() *tview.Flex {
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(newSummaryBox(), 0, 1, false)
}
