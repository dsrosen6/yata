package tui

import (
	"context"
	"fmt"

	"github.com/dsrosen6/yata/models"
)

func (a *App) addTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	return a.Repos.Tasks.Create(ctx, task)
}

func (a *App) updateTask(ctx context.Context, task *models.Task) error {
	if _, err := a.Repos.Tasks.Update(ctx, task); err != nil {
		return fmt.Errorf("updating task: %w", err)
	}

	return a.refreshTasks(ctx)
}

func (a *App) deleteTask(ctx context.Context, id int64) error {
	if err := a.Repos.Tasks.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}

	return a.refreshTasks(ctx)
}

func (a *App) refreshTasks(ctx context.Context) error {
	tasks, err := a.Repos.Tasks.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("listing tasks: %w", err)
	}

	a.Tasks = tasks
	if err := a.List.RefreshTasks(ctx); err != nil {
		return fmt.Errorf("refreshing list handler tasks: %w", err)
	}
	return nil
}
