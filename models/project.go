package models

import (
	"context"
	"time"
)

type Project struct {
	ID        int64
	Title     string
	ParentID  *int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProjectRepo interface {
	ListAll(ctx context.Context) ([]*Project, error)
	ListByParentID(ctx context.Context, parentID int64) ([]*Project, error)
	Get(ctx context.Context, id int64) (*Project, error)
	GetByTitle(ctx context.Context, title string) (*Project, error)
	Create(ctx context.Context, p *Project) (*Project, error)
	Update(ctx context.Context, p *Project) (*Project, error)
	Delete(ctx context.Context, id int64) error
}
