package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

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

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting user home directory: %w", err)
	}

	// TODO: change this to the standard spots depending on OS
	dir := filepath.Join(home, ".yata")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating app directory at %s: %w", dir, err)
	}

	d, err := sqlitedb.NewHandler(schema, filepath.Join(dir, "app.db"))
	if err != nil {
		return fmt.Errorf("initializing sqlite handler: %w", err)
	}

	defer func() {
		if err := d.Close(); err != nil {
			fmt.Println("Error closing handler:", err)
		}
	}()

	stores, err := d.InitStores(ctx)
	if err != nil {
		return fmt.Errorf("initializing repositories: %w", err)
	}

	return tui.Run(cfg, stores)
}
