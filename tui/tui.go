package tui

import (
	"context"
	"fmt"

	"github.com/dsrosen6/yata/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app *App

type App struct {
	*tview.Application
	MainFlex *tview.Flex
	ListFlex *tview.Flex
	Repos    *models.AllRepos
	List     *ListHandler
	Tasks    []*models.Task

	AddingTask bool
}

func init() {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
}

func NewApp(repos *models.AllRepos) *App {
	flx := tview.NewFlex()
	return &App{
		Application: tview.NewApplication().SetRoot(flx, true),
		MainFlex:    flx,
		Repos:       repos,
		Tasks:       []*models.Task{},
	}
}

func (a *App) Init(ctx context.Context) error {
	initialTasks, err := a.Repos.Tasks.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("initial task fetch: %w", err)
	}
	a.Tasks = initialTasks

	a.List = NewListHandler()
	if err := a.List.Init(ctx); err != nil {
		return fmt.Errorf("initializing list handler: %w", err)
	}

	a.ListFlex = newListFlex(a.List)
	a.MainFlex.AddItem(a.ListFlex, 0, 1, true).
		AddItem(newSummaryFlex(), 0, 1, false)

	a.SetFocus(a.ListFlex)
	return nil
}

func Run(ctx context.Context, repos *models.AllRepos) error {
	app = NewApp(repos)
	if err := app.Init(ctx); err != nil {
		return fmt.Errorf("initializing app: %w", err)
	}
	return app.Run()
}
