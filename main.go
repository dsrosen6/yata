package main

import (
	_ "embed"
	"fmt"

	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/sqlitedb"
	"github.com/dsrosen6/yata/tui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func run() error {
	h, err := sqlitedb.NewHandler(schema, "./app.db")
	if err != nil {
		return fmt.Errorf("initializing sqlite handler: %w", err)
	}

	defer func() {
		if err := h.Close(); err != nil {
			fmt.Println("Error closing handler:", err)
		}
	}()

	initialTasks := []*models.Task{
		{Title: "do something"},
		{Title: "do more things"},
		{Title: "do even more things"},
	}

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	flx := tview.NewFlex()
	list := tui.NewListHandler(initialTasks)
	box := tview.NewBox().SetTitle("summary").SetBorder(true)
	flx.AddItem(list, 0, 2, true)
	flx.AddItem(box, 0, 1, false)
	if err := tview.NewApplication().SetRoot(flx, true).Run(); err != nil {
		return fmt.Errorf("running app")
	}
	return nil
}
