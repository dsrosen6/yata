package tui

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/tui/models/form"
)

func projectsToItems(projects []*models.Project) []list.Item {
	// Start with the static "all" entry
	items := []list.Item{taskProjectItem{&models.Project{Title: "all"}}}
	for _, p := range projects {
		items = append(items, taskProjectItem{p})
	}
	return items
}

func (m *model) refreshProjects() tea.Cmd {
	return func() tea.Msg {
		currentIndex := m.projectList.Index()
		projects, err := m.stores.Projects.ListAll(context.Background())
		if err != nil {
			return storeErrorMsg{err}
		}
		items := append([]list.Item{}, projectsToItems(projects)...)
		cmd := m.projectList.SetItems(items)
		if currentIndex >= len(items) && len(items) > 0 {
			m.projectList.Select(len(items) - 1)
		}
		refreshedCmd := func() tea.Msg { return projectsRefreshedMsg{} }
		return tea.Sequence(cmd, refreshedCmd)
	}
}

func (m *model) insertProject(p taskProjectItem) tea.Cmd {
	return func() tea.Msg {
		if _, err := m.stores.Projects.Create(context.Background(), p.Project); err != nil {
			return storeErrorMsg{err}
		}

		return refreshProjectsMsg{}
	}
}

func (m *model) deleteProject(id int64) tea.Cmd {
	return func() tea.Msg {
		if err := m.stores.Projects.Delete(context.Background(), id); err != nil {
			return storeErrorMsg{err}
		}

		return refreshProjectsMsg{}
	}
}

func (m *model) selectedProject() taskProjectItem {
	item := m.projectList.SelectedItem()
	if item == nil {
		return taskProjectItem{}
	}
	return item.(taskProjectItem)
}

func (m *model) selectedProjectID() int64 {
	item := m.projectList.SelectedItem()
	if item == nil {
		return 0
	}
	sel := item.(taskProjectItem)
	if sel.Project == nil {
		return 0
	}

	return sel.ID
}

func newProjectEntryForm() (*form.Model, error) {
	fn := func(s string) error {
		if strings.TrimSpace(s) == "all" {
			return errors.New("'all' is a reserved filter and cannot be used")
		}
		return nil
	}

	fields := []form.Field{
		{
			Key:      "title",
			Required: true,
			Validate: fn,
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

func projectFromInputResult(r form.Result) taskProjectItem {
	t, ok := r["title"]
	if !ok {
		return taskProjectItem{}
	}

	return taskProjectItem{
		Project: &models.Project{
			Title: t,
		},
	}
}
