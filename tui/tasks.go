package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/tui/models/form"
)

func (m *model) refreshTasks() tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.stores.Tasks.ListAll(context.Background())
		if err != nil {
			return storeErrorMsg{err}
		}
		items := append([]list.Item{}, tasksToItems(tasks)...)
		return m.taskList.SetItems(items)
	}
}

func (m *model) insertTask(t taskItem) tea.Cmd {
	return func() tea.Msg {
		if _, err := m.stores.Tasks.Create(context.Background(), t.Task); err != nil {
			return storeErrorMsg{err}
		}

		return tea.BatchMsg{m.refreshTasks(), m.entryForm.Reset()}
	}
}

func (m *model) deleteTask(id int64) tea.Cmd {
	return func() tea.Msg {
		if err := m.stores.Tasks.Delete(context.Background(), id); err != nil {
			return storeErrorMsg{err}
		}

		return refreshTasksMsg{}
	}
}

func (m *model) toggleTaskComplete(t taskItem) tea.Cmd {
	return func() tea.Msg {
		t.Complete = !t.Complete
		if _, err := m.stores.Tasks.Update(context.Background(), t.Task); err != nil {
			return storeErrorMsg{err}
		}

		return refreshTasksMsg{}
	}
}

func (m *model) selectedTask() taskItem {
	return m.taskList.SelectedItem().(taskItem)
}

func (m *model) selectedTaskID() int64 {
	sel := m.taskList.SelectedItem().(taskItem)
	if sel.Task == nil {
		return 0
	}

	return sel.ID
}

func newTaskEntryForm(s styles) (*form.Model, error) {
	fields := []form.Field{
		{
			Key:      "title",
			Required: true,
		},
	}
	o := &form.Opts{
		Fields:           fields,
		PromptIfOneField: true,
		FocusedStyle:     s.focusedTextStyle,
		UnfocusedStyle:   s.unfocusedTextStyle,
		ErrorStyle:       s.errorTextStyle,
	}

	f, err := form.InitialInputModel(o)
	if err != nil {
		return nil, fmt.Errorf("creating model: %w", err)
	}

	return f, nil
}

func taskFromInputResult(r form.Result) taskItem {
	t, ok := r["title"]
	if !ok {
		return taskItem{}
	}

	return taskItem{
		Task: &models.Task{
			Title: t,
		},
	}
}

func tasksToItems(tasks []*models.Task) []list.Item {
	var items []list.Item
	for _, t := range tasks {
		items = append(items, taskItem{t})
	}
	return items
}
