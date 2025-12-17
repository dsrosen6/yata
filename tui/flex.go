package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *app) newListFlex(lh *tview.List) *tview.Flex {
	f := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(lh, 0, 2, true)

	f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			if !a.addingTask {
				a.addingTask = true
				f.AddItem(a.taskEntryField, 0, 1, true)
				a.SetFocus(a.taskEntryField)
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
