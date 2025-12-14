package tui

import (
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/dsrosen6/yata/models"
)

type (
	dbErrorMsg        struct{ error }
	refreshTasksMsg   string
	tasksRefreshedMsg struct{ tasks []*models.Task }
)

func createTaskForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title").
				Key("title").
				Validate(func(s string) error {
					if s == "" {
						return errors.New("title is required")
					}
					return nil
				}).
				Inline(true),
		),
	).WithShowHelp(false).WithTheme(huh.ThemeBase())
}

func (m *Model) insertTask(ctx context.Context, t *models.Task) tea.Cmd {
	return func() tea.Msg {
		if _, err := m.repos.Tasks.Create(ctx, t); err != nil {
			return dbErrorMsg{err}
		}

		return refreshTasksMsg("refresh")
	}
}

func (m *Model) toggleTaskComplete(ctx context.Context, t *models.Task) tea.Cmd {
	return func() tea.Msg {
		t.Complete = !t.Complete
		if _, err := m.repos.Tasks.Update(ctx, t); err != nil {
			return dbErrorMsg{err}
		}

		return refreshTasksMsg("refresh")
	}
}

func (m *Model) deleteTask(ctx context.Context, id int64) tea.Cmd {
	return func() tea.Msg {
		if err := m.repos.Tasks.Delete(ctx, id); err != nil {
			return dbErrorMsg{err}
		}

		return refreshTasksMsg("refresh")
	}
}

func (m *Model) refreshTasks(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.repos.Tasks.ListAll(ctx)
		if err != nil {
			return dbErrorMsg{err}
		}

		return tasksRefreshedMsg{tasks: tasks}
	}
}
