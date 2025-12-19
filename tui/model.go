package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	fbox "github.com/dsrosen6/tea-flexbox"
	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/tui/models/form"
)

type (
	model struct {
		stores *models.AllRepos
		tasks  []taskItem

		taskMode     taskMode
		entryForm    *form.Model
		sortParams   *models.SortParams
		taskList     list.Model
		listList     list.Model
		selected     map[int]struct{}
		currentFocus focus

		pendingAddTask bool
		dimensions
		styles
	}

	focus    int
	taskMode int
)

const (
	focusTasks focus = iota
	focusLists
	focusTaskEntry
	focusListEntry
)

func initialModel(s styles, stores *models.AllRepos) (*model, error) {
	entry, err := newTaskEntryForm(s)
	if err != nil {
		return nil, fmt.Errorf("creating task entry form: %w", err)
	}

	tasks, err := stores.Tasks.ListAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("getting initial tasks: %w", err)
	}

	return &model{
		styles:     s,
		taskList:   initialTaskList(tasks),
		stores:     stores,
		entryForm:  entry,
		sortParams: &models.SortParams{SortBy: models.SortByComplete},
		tasks:      []taskItem{},
		selected:   make(map[int]struct{}),
	}, nil
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(m.refreshTasks(), m.entryForm.Init())
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	switch m.currentFocus {
	case focusTasks:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "enter", " ":
				if len(m.taskList.Items()) > 0 {
					return m, m.toggleTaskComplete(m.selectedTask())
				}
			case "x":
				if len(m.taskList.Items()) > 0 {
					return m, m.deleteTask(m.selectedTaskID())
				}
			case "a":
				m.currentFocus = focusTaskEntry
				m.pendingAddTask = true
			}
		case refreshTasksMsg:
			return m, m.refreshTasks()
		case tasksRefreshedMsg:
			m.pendingAddTask = false
			return m, nil
		}

	case focusTaskEntry:
		f, cmd := m.entryForm.Update(msg)
		m.entryForm = f.(*form.Model)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.currentFocus = focusTasks
				m.pendingAddTask = false
				return m, m.entryForm.Reset()
			}
		case form.ResultMsg:
			m.currentFocus = focusTasks
			t := taskFromInputResult(msg.Result)
			return m, m.insertTask(t)
		}
		return m, cmd
	}

	var cmd tea.Cmd
	m.taskList, cmd = m.taskList.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	topBox := fbox.New(fbox.Horizontal, 1)
	topBox.AddTitleBox(m.createTasksBox(), 1, nil)

	fl := fbox.New(fbox.Vertical, 1).
		AddFlexBox(topBox, 7, nil).
		AddTitleBox(m.createEntryBox(), 1, func() bool { return m.currentFocus == focusTaskEntry })

	return fl.Render(m.width, m.height)
}
