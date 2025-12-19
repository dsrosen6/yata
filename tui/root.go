package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/config"
	"github.com/dsrosen6/yata/models"
)

var allStyles styles

type rootModel struct {
	cfg       *config.Config
	todoModel *model
	dimensions
}

type dimensions struct {
	width, height int
}

func Run(cfg *config.Config, stores *models.AllRepos) error {
	allStyles = generateStyles(cfg)
	m, err := newModel(cfg, stores)
	if err != nil {
		return fmt.Errorf("creating model: %w", err)
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return err
	}

	return nil
}

func newModel(cfg *config.Config, stores *models.AllRepos) (*rootModel, error) {
	td, err := initialModel(allStyles, stores)
	if err != nil {
		return nil, fmt.Errorf("creating todo list model: %w", err)
	}
	return &rootModel{
		cfg:       cfg,
		todoModel: td,
	}, nil
}

func (m *rootModel) Init() tea.Cmd {
	return m.todoModel.Init()
}

func (m *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var todoModel tea.Model
	todoModel, cmd = m.todoModel.Update(msg)
	m.todoModel = todoModel.(*model)

	return m, cmd
}

func (m *rootModel) View() string {
	return m.todoModel.View()
}
