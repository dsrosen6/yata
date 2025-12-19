package tui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	fbox "github.com/dsrosen6/tea-flexbox"
	"github.com/dsrosen6/tea-flexbox/titlebox"
	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/tui/models/form"
)

type (
	todoListModel struct {
		stores *models.AllRepos
		tasks  []*models.Task

		cursor    int
		taskMode  taskMode
		entryForm *form.Model
		selected  map[int]struct{}

		pendingAdd bool
		dimensions
		styles
	}

	storeErrorMsg     struct{ error }
	refreshTasksMsg   struct{}
	tasksRefreshedMsg struct{ tasks []*models.Task }

	taskMode int
)

const (
	taskModeViewing taskMode = iota
	taskModeCreating
)

func initialTodoList(s styles, stores *models.AllRepos) (*todoListModel, error) {
	entry, err := newTaskEntryForm(s)
	if err != nil {
		return nil, fmt.Errorf("creating task entry form: %w", err)
	}

	return &todoListModel{
		styles:    s,
		stores:    stores,
		entryForm: entry,
		tasks:     []*models.Task{},
		selected:  make(map[int]struct{}),
	}, nil
}

func (m *todoListModel) Init() tea.Cmd {
	return tea.Batch(m.refreshTasks(), m.entryForm.Init())
}

func (m *todoListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	switch m.taskMode {
	case taskModeViewing:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				m.cursor = cursorUp(m.cursor, len(m.tasks)-1)
			case "down", "j":
				m.cursor = cursorDown(m.cursor, len(m.tasks)-1)
			case "enter", " ":
				if len(m.tasks) > 0 {
					return m, m.toggleTaskComplete(m.tasks[m.cursor])
				}
			case "d":
				if len(m.tasks) > 0 {
					return m, m.deleteTask(m.tasks[m.cursor].ID)
				}
			case "a":
				m.taskMode = taskModeCreating
				m.pendingAdd = true
			}
		case refreshTasksMsg:
			return m, m.refreshTasks()
		case tasksRefreshedMsg:
			m.tasks = msg.tasks
			if len(m.tasks) == 0 {
				m.cursor = 0
			} else if m.pendingAdd {
				m.cursor = len(m.tasks) - 1
			} else if m.cursor >= len(m.tasks) {
				m.cursor = len(m.tasks) - 1
			}

			m.pendingAdd = false
			return m, nil
		}

	case taskModeCreating:
		f, cmd := m.entryForm.Update(msg)
		m.entryForm = f.(*form.Model)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.taskMode = taskModeViewing
				m.pendingAdd = false
				return m, m.entryForm.Reset()
			}
		case form.ResultMsg:
			m.taskMode = taskModeViewing
			t := taskFromInputResult(msg.Result)
			return m, m.insertTask(t)
		}
		return m, cmd
	}

	return m, nil
}

func (m *todoListModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	fl := fbox.New(fbox.Vertical, 1).
		AddTitleBox(m.createTasksBox(), 3, nil).
		AddTitleBox(m.createEntryBox(), 1, func() bool { return m.taskMode == taskModeCreating })

	return fl.Render(m.width, m.height)
}

func (m *todoListModel) createTasksBox() titlebox.Box {
	boxStyle := m.focusedBoxStyle
	titleStyle := m.focusedBoxTitleStyle
	if m.taskMode == taskModeCreating {
		boxStyle = m.unfocusedBoxStyle
		titleStyle = m.unfocusedBoxTitleStyle
	}

	return titlebox.New().
		SetTitle("tasks").
		SetBody(m.tasksOutput()).
		SetTitleAlignment(titlebox.AlignLeft).
		SetBoxStyle(boxStyle.Padding(0, 1)).
		SetTitleStyle(titleStyle)
}

func (m *todoListModel) createEntryBox() titlebox.Box {
	return titlebox.New().
		SetTitle("new task").
		SetBody(m.entryForm.View()).
		SetTitleAlignment(titlebox.AlignLeft).
		SetBoxStyle(m.focusedBoxStyle).
		SetTitleStyle(m.focusedBoxTitleStyle)
}

func (m *todoListModel) tasksOutput() string {
	if len(m.tasks) == 0 && m.taskMode != taskModeCreating && !m.pendingAdd {
		return "No tasks found\n"
	}

	uncheckedIcon := "󰄱"
	checkedIcon := "󰄵"
	var b strings.Builder
	for i, t := range m.tasks {
		checked := uncheckedIcon
		if t.Complete {
			checked = checkedIcon
		}

		s := fmt.Sprintf("%s %s", checked, t.Title)
		if m.cursor == i && m.taskMode != taskModeCreating {
			b.WriteString(m.focusedTextStyle.Render(s))
		} else {
			b.WriteString(m.unfocusedTextStyle.Render(s))
		}
		b.WriteRune('\n')
	}

	return b.String()
}

func (m *todoListModel) refreshTasks() tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.stores.Tasks.ListAll(context.Background())
		if err != nil {
			return storeErrorMsg{err}
		}

		return tasksRefreshedMsg{tasks: tasks}
	}
}

func (m *todoListModel) insertTask(t *models.Task) tea.Cmd {
	return func() tea.Msg {
		if _, err := m.stores.Tasks.Create(context.Background(), t); err != nil {
			return storeErrorMsg{err}
		}

		return tea.BatchMsg{m.refreshTasks(), m.entryForm.Reset()}
	}
}

func (m *todoListModel) deleteTask(id int64) tea.Cmd {
	return func() tea.Msg {
		if err := m.stores.Tasks.Delete(context.Background(), id); err != nil {
			return storeErrorMsg{err}
		}

		return refreshTasksMsg{}
	}
}

func (m *todoListModel) toggleTaskComplete(t *models.Task) tea.Cmd {
	return func() tea.Msg {
		t.Complete = !t.Complete
		if _, err := m.stores.Tasks.Update(context.Background(), t); err != nil {
			return storeErrorMsg{err}
		}

		return refreshTasksMsg{}
	}
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

func taskFromInputResult(r form.Result) *models.Task {
	t, ok := r["title"]
	if !ok {
		return nil
	}

	return &models.Task{
		Title: t,
	}
}
