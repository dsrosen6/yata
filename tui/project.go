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

type (
	refreshProjectsMsg    struct{ selectProjectID int64 }
	gotUpdatedProjectsMsg struct {
		projects        []list.Item
		selectProjectID int64
	}
)

func (m *model) checkProjectChanged() tea.Cmd {
	return func() tea.Msg {
		sel := m.selectedProjectID()
		if m.currentProjectID != sel {
			m.currentProjectID = sel
			return refreshTasksMsg{selectTaskID: 0}
		}
		return nil
	}
}

func (m *model) selectProject(id int64) tea.Cmd {
	for i, item := range m.projectList.Items() {
		if p, ok := item.(taskProjectItem); ok && p.ID == id {
			m.projectList.Select(i)
			break
		}
	}
	return nil
}

func (m *model) adjustProjectListIndex() tea.Cmd {
	currentIndex := m.projectList.Index()
	if currentIndex >= len(m.projectList.Items()) && len(m.projectList.Items()) > 0 {
		m.projectList.Select(len(m.projectList.Items()) - 1)
	}
	return nil
}

func (m *model) refreshProjects(selectProjectID int64) tea.Cmd {
	return func() tea.Msg {
		projects, err := m.stores.Projects.ListAll(context.Background())
		if err != nil {
			return storeErrorMsg{err}
		}
		items := append([]list.Item{}, projectsToItems(projects)...)
		return gotUpdatedProjectsMsg{
			projects:        items,
			selectProjectID: selectProjectID,
		}
	}
}

func (m *model) insertProject(p taskProjectItem) tea.Cmd {
	return func() tea.Msg {
		created, err := m.stores.Projects.Create(context.Background(), p.Project)
		if err != nil {
			return storeErrorMsg{err}
		}

		return refreshProjectsMsg{selectProjectID: created.ID}
	}
}

func (m *model) deleteProject(id int64) tea.Cmd {
	return func() tea.Msg {
		if err := m.stores.Projects.Delete(context.Background(), id); err != nil {
			return storeErrorMsg{err}
		}

		return refreshProjectsMsg{selectProjectID: 0}
	}
}

func (m *model) selectedProjectID() int64 {
	sel := m.selectedProject()
	if sel == nil {
		return 0
	}

	return sel.ID
}

func (m *model) selectedProject() *taskProjectItem {
	item := m.projectList.SelectedItem()
	if item == nil {
		return nil
	}
	sel := item.(taskProjectItem)
	if sel.Project == nil {
		return nil
	}

	return &sel
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

func projectsToItems(projects []*models.Project) []list.Item {
	// Start with the static "all" entry
	items := []list.Item{taskProjectItem{&models.Project{Title: "all"}}}
	for _, p := range projects {
		items = append(items, taskProjectItem{p})
	}
	return items
}
