package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/models"
)

type Model struct {
	repos *models.AllRepos
	tasks []*models.Task

	cursor   int
	selected map[int]struct{}
}

func InitialModel(r *models.AllRepos) *Model {
	return &Model{
		repos:    r,
		tasks:    []*models.Task{},
		selected: make(map[int]struct{}),
	}
}

func (m *Model) Init() tea.Cmd {
	// return m.fetchTasks(ctx)
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.cursor = cursorUp(m.cursor, len(m.tasks)-1)
		case "down", "j":
			m.cursor = cursorDown(m.cursor, len(m.tasks)-1)
		case "enter":
			m.tasks[m.cursor].Complete = !m.tasks[m.cursor].Complete
		}
	}

	return m, nil
}

func (m *Model) View() string {
	s := "Current Tasks:\n"

	for i, task := range m.tasks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if task.Complete {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, task.Title)
	}

	s += "\nctrl+c to quit"
	return s
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
