package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	fbox "github.com/dsrosen6/tea-flexbox"
	"github.com/dsrosen6/tea-flexbox/titlebox"
	"github.com/dsrosen6/yata/config"
	"github.com/dsrosen6/yata/models"
)

type model struct {
	cfg    *config.Config
	stores *models.AllRepos
	dimensions
	styles
}

type styles struct {
	mainStyle  lipgloss.Style
	titleStyle lipgloss.Style
}

type dimensions struct {
	width, height int
}

func Run(cfg *config.Config, stores *models.AllRepos) error {
	m := newModel(cfg, stores)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return err
	}

	return nil
}

func newModel(cfg *config.Config, stores *models.AllRepos) model {
	return model{
		cfg:    cfg,
		stores: stores,
		styles: styles{
			mainStyle:  mainStyle(cfg),
			titleStyle: boxTitleStyle(cfg),
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	return m.makeRootBox().Render(m.width, m.height)
}

func (m model) makeRootBox() fbox.Box {
	return fbox.NewBox(fbox.Horizontal, 1).
		AddItem(fbox.FlexBoxToItem(m.makeLeftBox(), 1)).
		AddItem(fbox.FlexBoxToItem(m.makeRightBox(), 1))
}

func (m model) makeLeftBox() fbox.Box {
	top := titlebox.New().
		SetTitle("tasks").
		SetBody("tasks will be here").
		SetTitleAlignment(titlebox.AlignLeft).
		SetBoxStyle(m.mainStyle).
		SetTitleStyle(m.titleStyle)

	return fbox.NewBox(fbox.Vertical, 1).
		AddItem(fbox.TitleBoxToItem(top, 1)).
		AddItem(fbox.StyleToItem(m.mainStyle, "placeholder", 1))
}

func (m model) makeRightBox() fbox.Box {
	return fbox.NewBox(fbox.Vertical, 1).
		AddItem(fbox.StyleToItem(m.mainStyle, "another placeholder", 1))
}
