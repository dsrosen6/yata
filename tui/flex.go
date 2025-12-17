package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *app) newListFlex(lh *tview.List) *tview.Flex {
	f := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(lh, 0, 1, true)

	f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			if !a.addingTask {
				a.addingTask = true
				a.rootFlex.AddItem(a.entryFlex, 0, 1, true)
				a.SetFocus(a.entryFlex)
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
