package main

import (
	"context"
	"database/sql"
	"fmt"
)

func initDB(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./app.db")
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	if _, err := db.ExecContext(ctx, schema); err != nil {
		return nil, fmt.Errorf("executing schema: %w", err)
	}

	return db, nil
}
