package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dsrosen6/yata/models"
	"github.com/dsrosen6/yata/tui/models/form"
)

func listsToItems(lists []*models.List) []list.Item {
	var items []list.Item
	for _, l := range lists {
		items = append(items, taskListItem{l})
	}
	return items
}

func (m *model) refreshLists() tea.Cmd {
	return func() tea.Msg {
		lists, err := m.stores.Lists.ListAll(context.Background())
		if err != nil {
			return storeErrorMsg{err}
		}
		items := append([]list.Item{}, listsToItems(lists)...)
		return m.listList.SetItems(items)
	}
}

func (m *model) insertList(l taskListItem) tea.Cmd {
	return func() tea.Msg {
		if _, err := m.stores.Lists.Create(context.Background(), l.List); err != nil {
			return storeErrorMsg{err}
		}

		return tea.BatchMsg{m.refreshLists(), m.listEntryForm.Reset()}
	}
}

func (m *model) deleteList(id int64) tea.Cmd {
	return func() tea.Msg {
		if err := m.stores.Lists.Delete(context.Background(), id); err != nil {
			return storeErrorMsg{err}
		}

		return refreshListsMsg{}
	}
}

func (m *model) selectedList() taskListItem {
	return m.listList.SelectedItem().(taskListItem)
}

func (m *model) selectedListID() int64 {
	sel := m.listList.SelectedItem().(taskListItem)
	if sel.List == nil {
		return 0
	}

	return sel.ID
}

func newListEntryForm() (*form.Model, error) {
	fields := []form.Field{
		{
			Key:      "title",
			Required: true,
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

func listFromInputResult(r form.Result) taskListItem {
	t, ok := r["title"]
	if !ok {
		return taskListItem{}
	}

	return taskListItem{
		List: &models.List{
			Title: t,
		},
	}
}
