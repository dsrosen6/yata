package tui

import "github.com/rivo/tview"

func newSummaryBox() *tview.Box {
	return tview.NewBox().
		SetTitle("summary").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
}
