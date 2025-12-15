package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/tui/input"
)

type (
	dbErrorMsg        struct{ error }
	refreshTasksMsg   string
	tasksRefreshedMsg struct{ tasks []*models.Task }
)

type taskEntryForm struct {
	Form *input.Model
}

func newTaskEntryForm(s styles) (*taskEntryForm, error) {
	o := &input.Opts{
		FieldKeys:        []string{"Title"},
		PromptIfOneField: false,
		FocusedStyle:     s.focusedStyle,
		UnfocusedStyle:   s.unfocusedStyle,
	}

	f, err := input.InitialInputModel(o)
	if err != nil {
		return nil, fmt.Errorf("creating model: %w", err)
	}

	return &taskEntryForm{
		Form: f,
	}, nil
}

func taskFromInputResult(r input.Result) *models.Task {
	t, ok := r["Title"]
	if !ok {
		return nil
	}

	return &models.Task{
		Title: t,
	}
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
