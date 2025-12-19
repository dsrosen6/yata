package models

import (
	"context"
	"time"
)

type List struct {
	ID        int64
	Title     string
	ParentID  *int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListRepo interface {
	ListAll(ctx context.Context) ([]*List, error)
	ListByParentID(ctx context.Context, parentID int64) ([]*List, error)
	Get(ctx context.Context, id int64) (*List, error)
	Create(ctx context.Context, l *List) (*List, error)
	Update(ctx context.Context, l *List) (*List, error)
	Delete(ctx context.Context, id int64) error
}
