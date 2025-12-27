package tui

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/models"
)

func (m *model) quit() tea.Cmd {
	return func() tea.Msg {
		if err := m.saveAppState(); err != nil {
			slog.Error("saving app state", "error", err)
		}
		return tea.Quit()
	}
}

func getAppState(stores *models.AllRepos) (*models.AppState, error) {
	ctx := context.Background()
	s, err := stores.AppState.Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Debug("no app state found; creating")
			s, err = stores.AppState.Create(ctx)
			if err != nil {
				return nil, fmt.Errorf("creating initial app state: %w", err)
			}
			slog.Debug("default app state created")
			return s, nil
		}

		return nil, fmt.Errorf("getting initial app state: %w", err)
	}

	return s, nil
}

func (m *model) saveAppState() error {
	ctx := context.Background()
	if _, err := m.stores.AppState.Update(ctx, m.state); err != nil {
		return fmt.Errorf("updating saved app state: %w", err)
	}
	slog.Debug("app state updated", logAppState(m.state))

	return nil
}

func logAppState(s *models.AppState) slog.Attr {
	var selProj int64
	if s.SelectedProjectID != nil {
		selProj = *s.SelectedProjectID
	}
	return slog.Group(
		"app state",
		slog.Int64("selected_project_id", selProj),
		slog.Bool("show_help", s.ShowHelp),
	)
}
