package tui

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	fbox "github.com/dsrosen6/tea-flexbox"
	"github.com/dsrosen6/tea-flexbox/titlebox"
	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/tui/models/form"
)

type (
	todoListModel struct {
		stores *models.AllRepos
		tasks  []todoItem

		cursor     int
		taskMode   taskMode
		entryForm  *form.Model
		sortParams *models.SortParams
		list       list.Model
		selected   map[int]struct{}

		pendingAdd bool
		dimensions
		styles
	}

	storeErrorMsg     struct{ error }
	refreshTasksMsg   struct{}
	tasksRefreshedMsg struct{}
	todoItem          struct{ *models.Task }
	todoItemDelegate  struct{}

	taskMode int
)

func (d todoItemDelegate) Height() int                             { return 1 }
func (d todoItemDelegate) Spacing() int                            { return 0 }
func (d todoItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d todoItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(todoItem)
	if !ok {
		return
	}

	checked := "󰄱"
	if i.Complete {
		checked = "󰄵"
	}

	str := fmt.Sprintf("%s %s", checked, i.Title)
	fn := allStyles.unfocusedTextStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return allStyles.focusedTextStyle.Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
func (t todoItem) FilterValue() string { return t.Title }

const (
	taskModeViewing taskMode = iota
	taskModeCreating
)

func initialTodoList(s styles, stores *models.AllRepos) (*todoListModel, error) {
	entry, err := newTaskEntryForm(s)
	if err != nil {
		return nil, fmt.Errorf("creating task entry form: %w", err)
	}

	tasks, err := stores.Tasks.ListAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("getting initial tasks: %w", err)
	}

	items := tasksToItems(tasks)
	l := list.New(items, todoItemDelegate{}, 10, 10)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	return &todoListModel{
		styles:     s,
		list:       l,
		stores:     stores,
		entryForm:  entry,
		sortParams: &models.SortParams{SortBy: models.SortByComplete},
		tasks:      []todoItem{},
		selected:   make(map[int]struct{}),
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
			case "enter", " ":
				if len(m.list.Items()) > 0 {
					return m, m.toggleTaskComplete(m.selectedTask())
				}
			case "x":
				if len(m.list.Items()) > 0 {
					return m, m.deleteTask(m.selectedTaskID())
				}
			case "a":
				m.taskMode = taskModeCreating
				m.pendingAdd = true
			}
		case refreshTasksMsg:
			return m, m.refreshTasks()
		case tasksRefreshedMsg:
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

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
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
		SetBody(m.list.View()).
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

func (m *todoListModel) refreshTasks() tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.stores.Tasks.ListAll(context.Background())
		if err != nil {
			return storeErrorMsg{err}
		}
		items := append([]list.Item{}, tasksToItems(tasks)...)
		return m.list.SetItems(items)
	}
}

func (m *todoListModel) insertTask(t todoItem) tea.Cmd {
	return func() tea.Msg {
		if _, err := m.stores.Tasks.Create(context.Background(), t.Task); err != nil {
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

func (m *todoListModel) toggleTaskComplete(t todoItem) tea.Cmd {
	return func() tea.Msg {
		t.Complete = !t.Complete
		if _, err := m.stores.Tasks.Update(context.Background(), t.Task); err != nil {
			return storeErrorMsg{err}
		}

		return refreshTasksMsg{}
	}
}

func (m *todoListModel) selectedTask() todoItem {
	return m.list.SelectedItem().(todoItem)
}

func (m *todoListModel) selectedTaskID() int64 {
	sel := m.list.SelectedItem().(todoItem)
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

func taskFromInputResult(r form.Result) todoItem {
	t, ok := r["title"]
	if !ok {
		return todoItem{}
	}

	return todoItem{
		Task: &models.Task{
			Title: t,
		},
	}
}

func tasksToItems(tasks []*models.Task) []list.Item {
	var items []list.Item
	for _, t := range tasks {
		items = append(items, todoItem{t})
	}
	return items
}
