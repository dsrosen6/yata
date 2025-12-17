package tui

import (
	"github.com/dsrosen6/yata/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func setInputColors(cfg *config.Config, f *tview.InputField) {
	f.SetFieldBackgroundColor(tcell.ColorDefault)
	f.SetLabelColor(cfg.MainColor)
	f.SetBorderColor(cfg.MainColor)
}

func setListColors(cfg *config.Config, l *tview.List) {
	l.SetSelectedBackgroundColor(tcell.ColorDefault)
	l.SetSelectedTextColor(cfg.MainColor)
}
