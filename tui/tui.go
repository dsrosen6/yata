package tui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/tui/input"
	"github.com/dsrosen6/yata/tui/style"
)

var ctx = context.Background()

type (
	Model struct {
		repos *models.AllRepos
		tasks []*models.Task
		size  size

		cursor     int
		taskMode   taskMode
		entryForm  *taskEntryForm
		selected   map[int]struct{}
		pendingAdd bool

		styles
	}

	styles struct {
		borderStyle    lipgloss.Style
		focusedStyle   lipgloss.Style
		unfocusedStyle lipgloss.Style
	}

	size struct {
		width  int
		height int
	}

	Opts struct {
		BorderColor    string
		FocusedColor   string
		UnfocusedColor string
	}

	taskMode string
)

const (
	taskModeViewing  taskMode = "viewing"
	taskModeCreating taskMode = "creating"
)

func InitialModel(r *models.AllRepos, opts Opts) (*Model, error) {
	return &Model{
		styles:   generateStyles(opts),
		repos:    r,
		tasks:    []*models.Task{},
		taskMode: taskModeViewing,
		selected: make(map[int]struct{}),
	}, nil
}

func generateStyles(o Opts) styles {
	return styles{
		borderStyle:    style.BorderStyle(o.BorderColor),
		focusedStyle:   style.FocusedStyle(o.FocusedColor),
		unfocusedStyle: style.UnfocusedStyle(o.UnfocusedColor),
	}
}

func (m *Model) Init() tea.Cmd {
	return m.refreshTasks(ctx)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Universal
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size.width = msg.Width
		m.size.height = msg.Height

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
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
			case "enter", " ":
				if len(m.tasks) > 0 {
					return m, m.toggleTaskComplete(ctx, m.tasks[m.cursor])
				}
			case "d":
				if len(m.tasks) > 0 {
					return m, m.deleteTask(ctx, m.tasks[m.cursor].ID)
				}
			case "a":
				m.taskMode = taskModeCreating
				m.pendingAdd = true
				m.entryForm, _ = newTaskEntryForm(m.styles)
				return m, m.entryForm.Form.Init()
			}

		case refreshTasksMsg:
			return m, m.refreshTasks(ctx)

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
		form, cmd := m.entryForm.Form.Update(msg)
		m.entryForm.Form = form.(*input.Model)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.taskMode = taskModeViewing
				m.entryForm, _ = newTaskEntryForm(m.styles)
				m.pendingAdd = false
				return m, nil
			}

		case input.ResultMsg:
			m.taskMode = taskModeViewing
			t := taskFromInputResult(msg.Result)
			return m, m.insertTask(ctx, t)
		}

		return m, cmd
	}

	return m, nil
}

func (m *Model) View() string {
	header := m.borderStyle.
		Align(lipgloss.Center).
		Width(m.size.width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Render("Header")

	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.size.width).
		Render("Footer")

	content := lipgloss.NewStyle().
		Width(m.size.width).
		Height(m.size.height-lipgloss.Height(header)-lipgloss.Height(footer)).
		Align(lipgloss.Center, lipgloss.Center).
		Render(m.tasksOutput())

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}

func (m *Model) tasksOutput() string {
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
			b.WriteString(m.focusedStyle.Render(s))
		} else {
			b.WriteString(m.unfocusedStyle.Render(s))
		}
		b.WriteRune('\n')
	}

	if m.taskMode == taskModeCreating {
		b.WriteString(m.focusedStyle.Render(uncheckedIcon + " "))
		b.WriteString(m.entryForm.Form.View())
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
