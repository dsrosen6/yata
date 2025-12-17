package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func setDefaultInputColors(f *tview.InputField) {
	f.SetFieldBackgroundColor(tcell.ColorDefault)
	f.SetLabelColor(tcell.ColorBlue)
}

func setDefaultListColors(l *tview.List) {
	l.SetSelectedBackgroundColor(tcell.ColorDefault)
	l.SetSelectedTextColor(tcell.ColorBlue)
}
