package tui

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/tui/models/form"
	fbox "github.com/dsrosen6/yata/tui/render/flexbox"
)

const (
	topBoxName    = "topBox"
	taskViewName  = "taskView"
	taskEntryName = "taskEntry"
	projViewName  = "projectView"
	projEntryName = "projectEntry"
	helpViewName  = "helpView"
)

type (
	model struct {
		state        *models.AppState
		stores       *models.AllRepos
		keys         keyMap
		help         help.Model
		taskList     list.Model
		projectList  list.Model
		sortParams   *models.SortParams
		currentFocus focus

		initFlags
		forms
		dimensions
	}

	initFlags struct {
		// gotSize guards project/task refreshes from occuring until
		// the initial WindowSizeMsg is received
		gotSize bool
	}

	forms struct {
		taskEntryForm    *form.Model
		projectEntryForm *form.Model
	}
	dimensionsCalculatedMsg struct{ dimensions }
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

	forms := forms{
		taskEntryForm:    te,
		projectEntryForm: pe,
	}

	s, err := getAppState(stores)
	if err != nil {
		return nil, fmt.Errorf("getting initial app state: %w", err)
	}
	slog.Debug("got initial app state", logAppState(s))

	return &model{
		state:       s,
		stores:      stores,
		keys:        defaultKeyMap,
		help:        help.New(),
		taskList:    initialTaskList(),
		projectList: initialProjectList(),
		sortParams:  &models.SortParams{SortBy: models.SortByComplete},
		forms:       forms,
	}, nil
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, m.calculateDimensions(msg.Width, msg.Height)
	case dimensionsCalculatedMsg:
		m.dimensions = msg.dimensions
		m.projectList.SetDelegate(projectItemDelegate{maxWidth: m.projDelegMaxW})
		m.projectList.SetHeight(m.listsH)
		m.taskList.SetHeight(m.listsH)
		m.logDimensions()
		if !m.gotSize {
			m.gotSize = true
			return m, m.refreshProjects(m.state.SelectedProjectID)
		}

	case changeFocusMsg:
		m.currentFocus = msg.focus
		return m, m.calculateDimensions(m.windowW, m.windowH)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			if !m.currentFocus.isEntry() {
				return m, m.quit()
			}
		case key.Matches(msg, m.keys.toggleHelp):
			if !m.currentFocus.isEntry() {
				m.state.ShowHelp = !m.state.ShowHelp
				return m, m.calculateDimensions(m.windowW, m.windowH)
			}
		case key.Matches(msg, m.keys.delete):
			switch m.currentFocus {
			case focusTasks:
				if len(m.taskList.Items()) > 0 {
					return m, m.deleteTask(m.selectedTaskID())
				}
			case focusProjects:
				// Don't allow deleting the "all" entry (which has ID 0)
				if len(m.projectList.Items()) > 0 && m.selectedProjectID() != nil {
					return m, m.deleteProject(*m.selectedProjectID())
				}
			}
		case key.Matches(msg, m.keys.focusProjects):
			if !m.currentFocus.isEntry() {
				return m, changeFocus(focusProjects)
			}
		case key.Matches(msg, m.keys.focusTasks):
			if !m.currentFocus.isEntry() {
				return m, changeFocus(focusTasks)
			}
		case key.Matches(msg, m.keys.newTask):
			if !m.currentFocus.isEntry() {
				return m, tea.Batch(m.taskEntryForm.Init(), changeFocus(focusTaskEntry))
			}
		case key.Matches(msg, m.keys.newProject):
			if !m.currentFocus.isEntry() {
				return m, tea.Batch(m.projectEntryForm.Init(), changeFocus(focusProjectEntry))
			}
		}
	case refreshTasksMsg:
		return m, m.getUpdatedTasks(m.state.SelectedProjectID, msg.selectTaskID)

	case gotUpdatedTasksMsg:
		cmds := []tea.Cmd{
			m.taskList.SetItems(msg.tasks),
			m.adjustTaskListIndex(),
			m.selectTask(msg.selectTaskID),
			m.calculateDimensions(m.windowW, m.windowH),
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

		if msg.selectProjectID != nil {
			cmds = append(cmds, m.selectProject(msg.selectProjectID))
		}

		return m, tea.Batch(cmds...)

	case selectedProjectChangedMsg:
		m.state.SelectedProjectID = msg.selected
		return m, m.getUpdatedTasks(m.state.SelectedProjectID, 0)
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
				return m, tea.Batch(m.taskEntryForm.Reset(), changeFocus(focusTasks))
			}

		case form.ResultMsg:
			t := taskFromInputResult(msg.Result)
			return m, tea.Batch(
				m.insertTask(t, m.state.SelectedProjectID),
				m.taskEntryForm.Reset(),
				changeFocus(focusTasks),
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
				return m, tea.Batch(m.projectEntryForm.Reset(), changeFocus(focusProjects))
			}

		case form.ResultMsg:
			p := projectFromInputResult(msg.Result)
			return m, tea.Batch(
				m.insertProject(p),
				m.projectEntryForm.Reset(),
				changeFocus(focusProjects),
			)
		}
		return m, cmd
	}

	return m, cmd
}

func (m *model) View() string {
	if !m.gotSize {
		return "Initializing..."
	}

	return m.createFlexbox().Render(m.windowW, m.windowH)
}

func (m *model) createFlexbox() *fbox.Box {
	hv := m.help.ShortHelpView(m.helpKeys())
	return fbox.New(fbox.Vertical, 1).
		AddFlexBox(m.createTopBox(), topBoxName, 7, nil, nil, nil).
		AddTitleBox(m.createTaskEntryBox(), taskEntryName, 1, nil, nil, func() bool { return m.currentFocus == focusTaskEntry }).
		AddTitleBox(m.createProjectEntryBox(), projEntryName, 1, nil, nil, func() bool { return m.currentFocus == focusProjectEntry }).
		AddStyleBox(helpStyle, helpViewName, hv, 1, nil, fbox.FixedSize(1), func() bool { return m.state.ShowHelp })
}
