package tui

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	fbox "github.com/dsrosen6/yata/tui/render/flexbox"
)

type dimensions struct {
	totalWidth          int
	totalHeight         int
	projectBoxWidth     int
	projectDelegateMaxW int
	listsHeight         int

	// all of the below are for debug purposes
	renderedWidth   int
	renderedHeight  int
	topBoxWidth     int
	topBoxHeight    int
	taskEntryWidth  int
	taskEntryHeight int
	projEntryWidth  int
	projEntryHeight int
	helpWidth       int
	helpHeight      int
}

func (m *model) calculateDimensions(w, h int) tea.Cmd {
	return func() tea.Msg {
		d := &dimensions{}
		d.totalWidth = w
		d.totalHeight = h

		box := m.createFlexbox()
		layouts := box.CalculateItemLayouts(w, h)
		lm := fbox.LayoutsToMap(layouts)

		if l, ok := lm[topBoxName]; ok {
			d.topBoxWidth = l.FullWidth
			d.topBoxHeight = l.FullHeight
		}

		if l, ok := lm[taskEntryName]; ok {
			d.taskEntryWidth = l.FullWidth
			d.taskEntryHeight = l.FullHeight
		}

		if l, ok := lm[projEntryName]; ok {
			d.projEntryWidth = l.FullWidth
			d.projEntryHeight = l.FullHeight
		}

		if l, ok := lm[helpViewName]; ok {
			d.helpWidth = l.FullWidth
			d.helpHeight = l.FullHeight
		}

		rendered := box.Render(w, h)
		d.renderedWidth = lipgloss.Width(rendered)
		d.renderedHeight = lipgloss.Height(rendered)

		tb := m.createTopBox()
		fw, fh := tb.GetAllItemsFrameSize()
		d.projectBoxWidth = 15
		d.projectDelegateMaxW = d.projectBoxWidth - fw
		d.listsHeight = d.topBoxHeight - fh
		return dimensionsCalculatedMsg{*d}
	}
}

func (m *model) logDimensions() {
	d := m.dimensions
	slog.Debug(
		"dimensions calculated",
		"focus", m.currentFocus.toString(),
		"total_w", d.totalWidth,
		"total_h", d.totalHeight,
		slog.Group("rendered", "width", d.renderedWidth, "height", d.renderedHeight),
		slog.Group("top_box", "width", d.topBoxWidth, "height", d.topBoxHeight),
		slog.Group("task_entry", "width", d.taskEntryWidth, "height", d.taskEntryHeight),
		slog.Group("project_entry", "width", d.projectBoxWidth, "height", d.projEntryHeight),
		slog.Group("help", "width", d.helpWidth, "height", d.helpHeight),
		"proj_box_w", d.projectBoxWidth,
		"proj_del_max_w", d.projectDelegateMaxW,
		"lists_h", d.listsHeight,
	)
}
