package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	fbox "github.com/dsrosen6/tea-flexbox"
	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/tui/models/form"
)

type (
	model struct {
		stores           *models.AllRepos
		keys             keyMap
		taskList         list.Model
		projectList      list.Model
		taskEntryForm    *form.Model
		projectEntryForm *form.Model
		sortParams       *models.SortParams
		currentFocus     focus
		currentProjectID int64

		dimensions
	}
)

func initialModel(stores *models.AllRepos) (*model, error) {
	te, err := newTaskEntryForm()
	if err != nil {
		return nil, fmt.Errorf("creating task entry form: %w", err)
	}

	pe, err := newProjectEntryForm()
	if err != nil {
		return nil, fmt.Errorf("creating project entry form: %w", err)
	}

	tasks, err := stores.Tasks.ListAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("getting initial tasks: %w", err)
	}

	projects, err := stores.Projects.ListAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("getting initial projects: %w", err)
	}

	return &model{
		stores:           stores,
		keys:             defaultKeyMap,
		taskList:         initialTaskList(tasks),
		projectList:      initialProjectList(projects),
		taskEntryForm:    te,
		projectEntryForm: pe,
		sortParams:       &models.SortParams{SortBy: models.SortByComplete},
	}, nil
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			if !m.currentFocus.isEntry() {
				return m, tea.Quit
			}
		case key.Matches(msg, m.keys.delete):
			switch m.currentFocus {
			case focusTasks:
				if len(m.taskList.Items()) > 0 {
					return m, m.deleteTask(m.selectedTaskID())
				}
			case focusProjects:
				// Don't allow deleting the "all" entry (which has ID 0)
				if len(m.projectList.Items()) > 0 && m.selectedProjectID() != 0 {
					return m, m.deleteProject(m.selectedProjectID())
				}
			}

		// focus-switching binds
		case key.Matches(msg, m.keys.focusProjects):
			if !m.currentFocus.isEntry() {
				m.currentFocus = focusProjects
			}
		case key.Matches(msg, m.keys.focusTasks):
			if !m.currentFocus.isEntry() {
				m.currentFocus = focusTasks
			}
		case key.Matches(msg, m.keys.newProject):
			if !m.currentFocus.isEntry() {
				m.currentFocus = focusProjectEntry
				return m, m.projectEntryForm.Init()
			}
		case key.Matches(msg, m.keys.newTask):
			if !m.currentFocus.isEntry() {
				m.currentFocus = focusTaskEntry
				return m, m.taskEntryForm.Init()
			}
		}
	case refreshTasksMsg:
		return m, m.getUpdatedTasks(m.currentProjectID)

	case gotUpdatedTasksMsg:
		return m, tea.Batch(
			m.taskList.SetItems(msg.tasks),
			m.adjustTaskListIndex(),
		)

	case refreshProjectsMsg:
		return m, m.refreshProjects()

	case gotUpdatedProjectsMsg:
		return m, tea.Batch(
			m.projectList.SetItems(msg.projects),
			m.checkProjectChanged(),
			m.adjustProjectListIndex(),
		)
	}

	switch m.currentFocus {
	case focusTasks:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.toggleTaskComplete):
				if len(m.taskList.Items()) > 0 {
					return m, m.toggleTaskComplete(m.selectedTask())
				}
			}
			m.taskList, cmd = m.taskList.Update(msg)
			return m, cmd
		}

	case focusProjects:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			m.projectList, cmd = m.projectList.Update(msg)
			return m, tea.Batch(cmd, m.checkProjectChanged())
		}

	case focusTaskEntry:
		f, cmd := m.taskEntryForm.Update(msg)
		m.taskEntryForm = f.(*form.Model)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.cancelEntry):
				m.currentFocus = focusTasks
				return m, m.taskEntryForm.Reset()
			}
		case form.ResultMsg:
			m.currentFocus = focusTasks
			t := taskFromInputResult(msg.Result)
			return m, tea.Batch(m.insertTask(t, m.currentProjectID), m.taskEntryForm.Reset())
		}
		return m, cmd

	case focusProjectEntry:
		f, cmd := m.projectEntryForm.Update(msg)
		m.projectEntryForm = f.(*form.Model)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.cancelEntry):
				m.currentFocus = focusProjects
				return m, m.projectEntryForm.Reset()
			}
		case form.ResultMsg:
			m.currentFocus = focusProjects
			p := projectFromInputResult(msg.Result)
			return m, m.insertProject(p)
		}
		return m, cmd
	}

	return m, cmd
}

func (m *model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	topBox := fbox.New(fbox.Horizontal, 4).
		AddTitleBox(m.createProjectsBox(), 1, nil).
		AddTitleBox(m.createTasksBox(), 10, nil)

	fl := fbox.New(fbox.Vertical, 1).
		AddFlexBox(topBox, 7, nil).
		AddTitleBox(m.createTaskEntryBox(), 1, func() bool { return m.currentFocus == focusTaskEntry }).
		AddTitleBox(m.createProjectEntryBox(), 1, func() bool { return m.currentFocus == focusProjectEntry })

	return fl.Render(m.width, m.height)
}
