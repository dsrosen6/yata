package main

import (
	"context"
	_ "embed"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/sqlitedb"
	"github.com/dsrosen6/yata/tui"
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
	ctx := context.Background()
	h, err := sqlitedb.NewHandler(schema, "./app.db")
	if err != nil {
		return fmt.Errorf("initializing sqlite handler: %w", err)
	}

	defer func() {
		if err := h.Close(); err != nil {
			fmt.Println("Error closing handler:", err)
		}
	}()

	r, err := h.InitStores(ctx)
	if err != nil {
		return fmt.Errorf("initializing stores: %w", err)
	}

	m := tui.InitialModel(r)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("running tui: %w", err)
	}
	return nil
}
