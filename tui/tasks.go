package tui

import (
	"context"
	"fmt"

	"github.com/dsrosen6/yata/models"
)

func (a *app) addTask(ctx context.Context, task *models.Task) error {
	if _, err := a.repos.Tasks.Create(ctx, task); err != nil {
		return fmt.Errorf("creating task: %w", err)
	}

	return a.refreshTasks(ctx)
}

func (a *app) updateTask(ctx context.Context, task *models.Task) error {
	if _, err := a.repos.Tasks.Update(ctx, task); err != nil {
		return fmt.Errorf("updating task: %w", err)
	}

	return a.refreshTasks(ctx)
}

func (a *app) deleteTask(ctx context.Context, id int64) error {
	if err := a.repos.Tasks.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}

	return a.refreshTasks(ctx)
}

func (a *app) refreshTasks(ctx context.Context) error {
	tasks, err := a.repos.Tasks.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("listing tasks: %w", err)
	}

	a.tasks = tasks
	if err := a.refreshListTasks(a.taskList); err != nil {
		return fmt.Errorf("refreshing list handler tasks: %w", err)
	}
	return nil
}
