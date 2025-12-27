package sqlitedb

import (
	"context"
	"errors"

	"github.com/dsrosen6/yata/models"
)

type AppStateRepo struct {
	q *Queries
}

func NewAppStateRepo(q *Queries) *AppStateRepo {
	return &AppStateRepo{
		q: q,
	}
}

func (ar *AppStateRepo) Get(ctx context.Context) (*models.AppState, error) {
	ds, err := ar.q.GetAppState(ctx)
	if err != nil {
		return nil, err
	}

	return dbStateToState(ds), nil
}

func (ar *AppStateRepo) Create(ctx context.Context) (*models.AppState, error) {
	ds, err := ar.q.CreateAppState(ctx)
	if err != nil {
		return nil, err
	}

	return dbStateToState(ds), nil
}

func (ar *AppStateRepo) Update(ctx context.Context, s *models.AppState) (*models.AppState, error) {
	if s == nil {
		return nil, errors.New("received nil app state")
	}
	ds, err := ar.q.UpdateAppState(ctx, stateToUpdateParams(s))
	if err != nil {
		return nil, err
	}

	return dbStateToState(ds), nil
}

func stateToUpdateParams(s *models.AppState) *UpdateAppStateParams {
	return &UpdateAppStateParams{
		SelectedProjectID: s.SelectedProjectID,
		ShowHelp:          s.ShowHelp,
	}
}

func dbStateToState(s *AppState) *models.AppState {
	return &models.AppState{
		ID:                s.ID,
		SelectedProjectID: s.SelectedProjectID,
		ShowHelp:          s.ShowHelp,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}
