package main

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/dsrosen6/yata/config"
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
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	d, err := sqlitedb.NewHandler(schema, "./app.db")
	if err != nil {
		return fmt.Errorf("initializing sqlite handler: %w", err)
	}

	defer func() {
		if err := d.Close(); err != nil {
			fmt.Println("Error closing handler:", err)
		}
	}()

	repos, err := d.InitStores(ctx)
	if err != nil {
		return fmt.Errorf("initializing repositories: %w", err)
	}

	return tui.Run(ctx, cfg, repos)
}
