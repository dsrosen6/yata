package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	fbox "github.com/dsrosen6/tea-flexbox"
	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/tui/models/form"
)

type (
	model struct {
		stores *models.AllRepos

		taskList      list.Model
		listList      list.Model
		taskEntryForm *form.Model
		listEntryForm *form.Model
		sortParams    *models.SortParams
		currentFocus  focus

		dimensions
	}

	focus int
)

const (
	focusTasks focus = iota
	focusLists
	focusTaskEntry
	focusListEntry
)

func initialModel(stores *models.AllRepos) (*model, error) {
	te, err := newTaskEntryForm()
	if err != nil {
		return nil, fmt.Errorf("creating task entry form: %w", err)
	}

	le, err := newListEntryForm()
	if err != nil {
		return nil, fmt.Errorf("creating list entry form: %w", err)
	}

	tasks, err := stores.Tasks.ListAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("getting initial tasks: %w", err)
	}

	lists, err := stores.Lists.ListAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("getting initial lists: %w", err)
	}

	return &model{
		stores:        stores,
		taskList:      initialTaskList(tasks),
		listList:      initialListList(lists),
		taskEntryForm: te,
		listEntryForm: le,
		sortParams:    &models.SortParams{SortBy: models.SortByComplete},
	}, nil
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(m.refreshTasks(), m.taskEntryForm.Init())
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if m.currentFocus == focusTasks || m.currentFocus == focusLists {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
	}

	switch m.currentFocus {
	case focusTasks:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter", " ":
				if len(m.taskList.Items()) > 0 {
					return m, m.toggleTaskComplete(m.selectedTask())
				}
			case "x":
				if len(m.taskList.Items()) > 0 {
					return m, m.deleteTask(m.selectedTaskID())
				}
			case "a":
				m.currentFocus = focusTaskEntry
				return m, m.taskEntryForm.Init()
			case "1":
				m.currentFocus = focusLists
			}
		case refreshTasksMsg:
			return m, m.refreshTasks()
		}
	case focusLists:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "x":
				if len(m.listList.Items()) > 0 {
					return m, m.deleteList(m.selectedListID())
				}
			case "a":
				m.currentFocus = focusListEntry
				return m, m.listEntryForm.Init()
			case "2":
				m.currentFocus = focusTasks
			}
		case refreshListsMsg:
			return m, m.refreshLists()
		}

	case focusTaskEntry:
		f, cmd := m.taskEntryForm.Update(msg)
		m.taskEntryForm = f.(*form.Model)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.currentFocus = focusTasks
				return m, m.taskEntryForm.Reset()
			}
		case form.ResultMsg:
			m.currentFocus = focusTasks
			t := taskFromInputResult(msg.Result)
			return m, m.insertTask(t)
		}
		return m, cmd

	case focusListEntry:
		f, cmd := m.listEntryForm.Update(msg)
		m.listEntryForm = f.(*form.Model)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.currentFocus = focusLists
				return m, m.listEntryForm.Reset()
			}
		case form.ResultMsg:
			m.currentFocus = focusLists
			l := listFromInputResult(msg.Result)
			return m, m.insertList(l)
		}
		return m, cmd
	}

	var tlCmd, llCmd tea.Cmd
	m.taskList, tlCmd = m.taskList.Update(msg)
	m.listList, llCmd = m.listList.Update(msg)
	return m, tea.Batch(tlCmd, llCmd)
}

func (m *model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	topBox := fbox.New(fbox.Horizontal, 4).
		AddTitleBox(m.createListsBox(), 1, nil).
		AddTitleBox(m.createTasksBox(), 10, nil)

	fl := fbox.New(fbox.Vertical, 1).
		AddFlexBox(topBox, 7, nil).
		AddTitleBox(m.createTaskEntryBox(), 1, func() bool { return m.currentFocus == focusTaskEntry }).
		AddTitleBox(m.createListEntryBox(), 1, func() bool { return m.currentFocus == focusListEntry })

	return fl.Render(m.width, m.height)
}
