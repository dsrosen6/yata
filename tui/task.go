package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/tui/models/form"
)

type (
	refreshTasksMsg struct {
		projectID    int64
		selectTaskID int64
	}
	gotUpdatedTasksMsg struct {
		tasks        []list.Item
		selectTaskID int64
	}
)

// selectTask selects the task in the list with the provided ID. If no ID is provided,
// it takes the 0 index.
func (m *model) selectTask(id int64) tea.Cmd {
	if id == 0 {
		m.taskList.Select(0)
	}

	for i, item := range m.taskList.Items() {
		if t, ok := item.(taskItem); ok && t.ID == id {
			m.taskList.Select(i)
			break
		}
	}
	return nil
}

// adjustTaskListIndex is a safeguard for when the selected index of a task list becomes
// higher than total items. For example, when you delete the last item in the list, it shifts
// the selected index to the new last item in the list.
func (m *model) adjustTaskListIndex() tea.Cmd {
	currentIndex := m.taskList.Index()
	if currentIndex >= len(m.taskList.Items()) && len(m.taskList.Items()) > 0 {
		m.taskList.Select(len(m.taskList.Items()) - 1)
	}
	return nil
}

// getUpdatedTasks retrieves tasks from the store and returns a gotUpdatedTasksMsg with those tasks,
// which will later be used to update the visible list. It takes a project ID to filter by project,
// and a task ID to be later used to select the proper task once refreshed.
func (m *model) getUpdatedTasks(projectID, selectTaskID int64) tea.Cmd {
	return func() tea.Msg {
		var (
			tasks []*models.Task
			err   error
		)

		ctx := context.Background()

		// if no project ID is provided, assume the "all" view in the projects list is selected
		if m.currentProjectID == 0 {
			tasks, err = m.stores.Tasks.ListAll(ctx)
		} else {
			tasks, err = m.stores.Tasks.ListByProjectID(ctx, projectID)
		}

		if err != nil {
			return storeErrorMsg{err}
		}

		items := append([]list.Item{}, tasksToItems(tasks)...)
		return gotUpdatedTasksMsg{
			tasks:        items,
			selectTaskID: selectTaskID,
		}
	}
}

func (m *model) insertTask(t taskItem, projectID int64) tea.Cmd {
	return func() tea.Msg {
		if projectID != 0 {
			t.ProjectID = &projectID
		}

		created, err := m.stores.Tasks.Create(context.Background(), t.Task)
		if err != nil {
			return storeErrorMsg{err}
		}

		return refreshTasksMsg{
			projectID:    projectID,
			selectTaskID: created.ID,
		}
	}
}

func (m *model) deleteTask(id int64) tea.Cmd {
	return func() tea.Msg {
		if err := m.stores.Tasks.Delete(context.Background(), id); err != nil {
			return storeErrorMsg{err}
		}

		return refreshTasksMsg{selectTaskID: 0}
	}
}

func (m *model) toggleTaskComplete(t taskItem) tea.Cmd {
	return func() tea.Msg {
		t.Complete = !t.Complete
		if _, err := m.stores.Tasks.Update(context.Background(), t.Task); err != nil {
			return storeErrorMsg{err}
		}

		return refreshTasksMsg{selectTaskID: 0}
	}
}

func (m *model) selectedTask() taskItem {
	item := m.taskList.SelectedItem()
	if item == nil {
		return taskItem{}
	}
	return item.(taskItem)
}

func (m *model) selectedTaskID() int64 {
	item := m.taskList.SelectedItem()
	if item == nil {
		return 0
	}
	sel := item.(taskItem)
	if sel.Task == nil {
		return 0
	}

	return sel.ID
}

func newTaskEntryForm() (*form.Model, error) {
	fields := []form.Field{
		{
			Key:      "title",
			Required: true,
		},
	}
	o := &form.Opts{
		Fields:           fields,
		PromptIfOneField: true,
		FocusedStyle:     allStyles.focusedTextStyle,
		UnfocusedStyle:   allStyles.unfocusedTextStyle,
		ErrorStyle:       allStyles.errorTextStyle,
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
