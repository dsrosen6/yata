package main

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/dsrosen6/yata/sqlitedb"
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

	db, err := initDB(ctx)
	if err != nil {
		return fmt.Errorf("initializing db: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing db:", err)
		}
	}()

	q := sqlitedb.New(db)
	_ = sqlitedb.NewRepos(q)
	return nil
}
