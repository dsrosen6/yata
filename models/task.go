package models

import (
	"context"
	"time"
)

type Task struct {
	ID           int64
	Title        string
	ParentTaskID *int64
	ListID       *int64
	Complete     bool
	DueAt        *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TaskRepo interface {
	ListAll(ctx context.Context) ([]*Task, error)
	ListByListID(ctx context.Context, listID int64) ([]*Task, error)
	ListByParentID(ctx context.Context, parentID int64) ([]*Task, error)
	Get(ctx context.Context, id int64) (*Task, error)
	Create(ctx context.Context, t *Task) (*Task, error)
	Update(ctx context.Context, t *Task) (*Task, error)
	Delete(ctx context.Context, id int64) error
}
