package tui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/dsrosen6/yata/models"
)

var ctx = context.Background()

type Model struct {
	debug bool
	repos *models.AllRepos
	tasks []*models.Task

	cursor    int
	taskMode  taskMode
	entryForm *huh.Form
	selected  map[int]struct{}
}

type taskMode string

const (
	taskModeViewing  taskMode = "viewing"
	taskModeCreating taskMode = "creating"
)

func InitialModel(r *models.AllRepos) *Model {
	return &Model{
		repos:     r,
		tasks:     []*models.Task{},
		taskMode:  taskModeViewing,
		entryForm: createTaskForm(),
		selected:  make(map[int]struct{}),
	}
}

func (m *Model) Init() tea.Cmd {
	return m.refreshTasks(ctx)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Universal quit bind
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	switch m.taskMode {
	case taskModeViewing:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				m.cursor = cursorUp(m.cursor, len(m.tasks)-1)
			case "down", "j":
				m.cursor = cursorDown(m.cursor, len(m.tasks)-1)
			case "enter":
				if len(m.tasks) > 0 {
					return m, m.toggleTaskComplete(ctx, m.tasks[m.cursor])
				}
			case "d":
				if len(m.tasks) > 0 {
					return m, m.deleteTask(ctx, m.tasks[m.cursor].ID)
				}
			case "a":
				m.taskMode = taskModeCreating
				m.entryForm = createTaskForm()
				return m, m.entryForm.Init()
			}

		case refreshTasksMsg:
			return m, m.refreshTasks(ctx)

		case tasksRefreshedMsg:
			m.tasks = msg.tasks
			return m, nil
		}

	case taskModeCreating:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.taskMode = taskModeViewing
				return m, nil
			}
		}

		if m.entryForm != nil {
			form, cmd := m.entryForm.Update(msg)
			if f, ok := form.(*huh.Form); ok {
				m.entryForm = f
			}

			if m.entryForm.State == huh.StateCompleted {
				t := &models.Task{
					Title: m.entryForm.GetString("title"),
				}

				m.tasks = append(m.tasks, t)
				m.taskMode = taskModeViewing
				return m, m.insertTask(ctx, t)
			}
			return m, cmd
		}
	}

	return m, nil
}

func (m *Model) View() string {
	var b strings.Builder
	b.WriteString(tasksOutput(m.cursor, m.tasks))

	if m.taskMode == taskModeCreating {
		fmt.Fprintf(&b, "%s\n", m.entryForm.View())
	}

	if m.debug {
		fmt.Fprintf(&b, "Task Mode: %s\nForm State: %v\n", m.taskMode, m.entryForm.State)
	}

	return b.String()
}

func tasksOutput(cursor int, tasks []*models.Task) string {
	if len(tasks) == 0 {
		return "No tasks found\n"
	}

	var b strings.Builder
	b.WriteString("Current Tasks:\n")
	for i, t := range tasks {
		cstr := " "
		if cursor == i {
			cstr = ">"
		}

		checked := "󰄱"
		if t.Complete {
			checked = "󰄵"
		}
		s := fmt.Sprintf("%s %s %s\n", cstr, checked, t.Title)
		b.WriteString(s)
	}

	return b.String()
}

func cursorUp(c, top int) int {
	if c > 0 {
		c--
	} else {
		c = top
	}

	return c
}

func cursorDown(c, top int) int {
	if c < top {
		c++
	} else {
		c = 0
	}
	return c
}
