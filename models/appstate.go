package models

import (
	"context"
	"time"
)

type AppState struct {
	ID                int64
	SelectedProjectID *int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type AppStateRepo interface {
	Get(ctx context.Context) (*AppState, error)
	Create(ctx context.Context) (*AppState, error)
	Update(ctx context.Context, s *AppState) (*AppState, error)
}
