package sqlitedb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dsrosen6/yata/models"
	_ "modernc.org/sqlite"
)

type Handler struct {
	embedSchema string
	db          *sql.DB
	queries     *Queries
}

func NewHandler(embedSchema, dbPath string) (*Handler, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	q := New(db)
	return &Handler{
		embedSchema: embedSchema,
		db:          db,
		queries:     q,
	}, nil
}

func (h *Handler) InitStores(ctx context.Context) (*models.AllRepos, error) {
	if _, err := h.db.ExecContext(ctx, h.embedSchema); err != nil {
		return nil, fmt.Errorf("executing schema: %w", err)
	}

	return NewRepos(h.queries), nil
}

func (h *Handler) Close() error {
	return h.db.Close()
}

func NewRepos(q *Queries) *models.AllRepos {
	return &models.AllRepos{
		Tasks: NewTaskRepo(q),
		Lists: NewListRepo(q),
	}
}
