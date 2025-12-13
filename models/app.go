package models

import "context"

type StoreHandler interface {
	CreateRepos(ctx context.Context) (*AllRepos, error)
	Close() error
}

type AllRepos struct {
	Tasks TaskRepo
}
