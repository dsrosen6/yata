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
	rootFlex       *tview.Flex
	mainFlex       *tview.Flex
	entryFlex      *tview.Flex
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
	a.rootFlex = tview.NewFlex().SetDirection(tview.FlexRow)
	a.mainFlex = tview.NewFlex()
	a.taskList = a.newTaskList()
	a.listFlex = a.newListFlex(a.taskList)
	a.taskEntryField = a.newTaskEntryField()
	a.entryFlex = tview.NewFlex().AddItem(a.taskEntryField, 0, 1, true)
	a.repos = repos
	a.rootFlex.AddItem(a.mainFlex, 0, 8, true)
	a.rootFlex.SetInputCapture(a.globalInputCapture)
	a.Application = tview.NewApplication().SetRoot(a.rootFlex, true)

	return a
}

func (a *app) globalInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q':
		if !a.addingTask {
			a.Stop()
		}
	}
	return event
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
