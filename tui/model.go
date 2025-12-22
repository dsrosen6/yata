package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/tui/models/form"
	fbox "github.com/dsrosen6/yata/tui/render/flexbox"
)

type (
	model struct {
		stores           *models.AllRepos
		keys             keyMap
		help             help.Model
		showHelp         bool
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

type dimensions struct {
	totalWidth          int
	totalHeight         int
	projectBoxWidth     int
	projectDelegateMaxW int
}

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
		help:             help.New(),
		showHelp:         true,
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

func (m *model) calculateDimensions(msg tea.WindowSizeMsg) dimensions {
	d := &dimensions{}
	d.totalWidth = msg.Width
	d.totalHeight = msg.Height
	f, _ := m.createProjectsBox().FrameSize()

	d.projectBoxWidth = 15
	d.projectDelegateMaxW = d.projectBoxWidth - f
	return *d
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.dimensions = m.calculateDimensions(msg)
		m.projectList.SetDelegate(projectItemDelegate{maxWidth: m.projectDelegateMaxW})

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			if !m.currentFocus.isEntry() {
				return m, tea.Quit
			}
		case key.Matches(msg, m.keys.toggleHelp):
			if !m.currentFocus.isEntry() {
				m.showHelp = !m.showHelp
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
		case key.Matches(msg, m.keys.focusProjects):
			if m.currentFocus == focusTasks {
				m.currentFocus = focusProjects
			}
		case key.Matches(msg, m.keys.focusTasks):
			if m.currentFocus == focusProjects {
				m.currentFocus = focusTasks
			}
		case key.Matches(msg, m.keys.focusProjects):
			if !m.currentFocus.isEntry() {
				m.currentFocus = focusProjects
			}
		case key.Matches(msg, m.keys.focusTasks):
			if !m.currentFocus.isEntry() {
				m.currentFocus = focusTasks
			}
		case key.Matches(msg, m.keys.newItem):
			switch m.currentFocus {
			case focusProjects:
				m.currentFocus = focusProjectEntry
				return m, m.projectEntryForm.Init()
			case focusTasks:
				m.currentFocus = focusTaskEntry
				return m, m.taskEntryForm.Init()
			}
		}
	case refreshTasksMsg:
		return m, m.getUpdatedTasks(m.currentProjectID, msg.selectTaskID)

	case gotUpdatedTasksMsg:
		cmds := []tea.Cmd{
			m.taskList.SetItems(msg.tasks),
			m.adjustTaskListIndex(),
		}

		if msg.selectTaskID != 0 {
			cmds = append(cmds, m.selectTask(msg.selectTaskID))
		}

		return m, tea.Batch(cmds...)

	case refreshProjectsMsg:
		// Some commands that refresh the projects list will provide a
		// selected project ID. This is for cases like adding a project.
		// The project ID is passed down the command chain and then it is
		// set as the selected project once it reaches the end.
		return m, m.refreshProjects(msg.selectProjectID)

	case gotUpdatedProjectsMsg:
		cmds := []tea.Cmd{
			m.projectList.SetItems(msg.projects),
			m.checkProjectChanged(),
			m.adjustProjectListIndex(),
		}

		if msg.selectProjectID != 0 {
			cmds = append(cmds, m.selectProject(msg.selectProjectID))
		}

		return m, tea.Batch(cmds...)
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
			return m, tea.Batch(
				m.insertTask(t, m.currentProjectID),
				m.taskEntryForm.Reset(),
			)
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
			return m, tea.Batch(
				m.insertProject(p),
				m.projectEntryForm.Reset(),
			)
		}
		return m, cmd
	}

	return m, cmd
}

func (m *model) View() string {
	if m.totalWidth == 0 || m.totalHeight == 0 {
		return "Initializing..."
	}

	topBox := fbox.New(fbox.Horizontal, 4).
		AddTitleBox(m.createProjectsBox(), 1, fbox.FixedSize(m.projectBoxWidth), nil, nil).
		AddTitleBox(m.createTasksBox(), 8, nil, nil, nil)

	hv := m.help.ShortHelpView(m.helpKeys())
	fl := fbox.New(fbox.Vertical, 1).
		AddFlexBox(topBox, 7, nil, nil, nil).
		AddTitleBox(m.createTaskEntryBox(), 1, nil, nil, func() bool { return m.currentFocus == focusTaskEntry }).
		AddTitleBox(m.createProjectEntryBox(), 1, nil, nil, func() bool { return m.currentFocus == focusProjectEntry }).
		AddStyleBox(helpStyle, hv, 1, nil, fbox.FixedSize(1), func() bool { return m.showHelp })

	return fl.Render(m.totalWidth, m.totalHeight)
}
