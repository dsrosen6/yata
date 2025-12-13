package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/models"
)

type (
	gotTasksMsg []*models.Task
	storeErrMsg struct{ err error }
)

func (m *Model) fetchTasks(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.repos.Tasks.ListAll(ctx)
		if err != nil {
			return storeErrMsg{err}
		}

		return gotTasksMsg(tasks)
	}
}
