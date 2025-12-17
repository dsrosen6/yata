package tui

import (
	"context"
	"fmt"

	"github.com/dsrosen6/yata/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type app struct {
	*tview.Application
	repos          *models.AllRepos
	mainFlex       *tview.Flex
	listFlex       *tview.Flex
	taskList       *tview.List
	taskEntryField *tview.InputField
	tasks          []*models.Task

	addingTask bool
}

func init() {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
}

func Run(ctx context.Context, repos *models.AllRepos) error {
	a := newApp(repos)
	if err := a.init(ctx); err != nil {
		return fmt.Errorf("initializing app: %w", err)
	}
	return a.Run()
}

func newApp(repos *models.AllRepos) *app {
	a := &app{}
	a.mainFlex = tview.NewFlex()
	a.taskList = a.newTaskList()
	a.listFlex = a.newListFlex(a.taskList)
	a.taskEntryField = a.newTaskEntryField()
	a.repos = repos
	a.Application = tview.NewApplication().SetRoot(a.mainFlex, true)

	return a
}

func (a *app) init(ctx context.Context) error {
	tasks, err := a.repos.Tasks.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("initial task fetch: %w", err)
	}
	a.tasks = tasks

	// provides an initial load of tasks into the list
	if err := a.initTaskList(a.taskList); err != nil {
		return fmt.Errorf("initializing list handler: %w", err)
	}

	a.mainFlex.AddItem(a.listFlex, 0, 1, true).
		AddItem(newSummaryFlex(), 0, 1, false)

	a.SetFocus(a.listFlex)
	return nil
}
