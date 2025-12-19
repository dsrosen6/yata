package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/config"
	"github.com/dsrosen6/yata/models"
)

type model struct {
	cfg       *config.Config
	stores    *models.AllRepos
	todoModel *todoListModel
	dimensions
	styles
}

type dimensions struct {
	width, height int
}

func Run(cfg *config.Config, stores *models.AllRepos) error {
	m, err := newModel(cfg, stores)
	if err != nil {
		return fmt.Errorf("creating model: %w", err)
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return err
	}

	return nil
}

func newModel(cfg *config.Config, stores *models.AllRepos) (*model, error) {
	s := generateStyles(cfg)
	td, err := initialTodoList(s, stores)
	if err != nil {
		return nil, fmt.Errorf("creating todo list model: %w", err)
	}
	return &model{
		cfg:       cfg,
		stores:    stores,
		todoModel: td,
		styles:    s,
	}, nil
}

func (m *model) Init() tea.Cmd {
	return m.todoModel.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var todoModel tea.Model
	todoModel, cmd = m.todoModel.Update(msg)
	m.todoModel = todoModel.(*todoListModel)

	return m, cmd
}

func (m *model) View() string {
	return m.todoModel.View()
}
