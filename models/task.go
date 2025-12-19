package models

import (
	"context"
	"time"
)

type Task struct {
	ID        int64
	Title     string
	Complete  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TaskRepo interface {
	ListAll(ctx context.Context) ([]*Task, error)
	Get(ctx context.Context, id int64) (*Task, error)
	Create(ctx context.Context, t *Task) (*Task, error)
	Update(ctx context.Context, t *Task) (*Task, error)
	Delete(ctx context.Context, id int64) error
}
