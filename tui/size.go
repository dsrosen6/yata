package tui

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	fbox "github.com/dsrosen6/yata/tui/render/flexbox"
)

type dimensions struct {
	windowW       int
	windowH       int
	projBoxW      int
	projDelegMaxW int
	listsH        int

	// all of the below are for debug purposes
	renderedW       int
	renderedH       int
	topBoxMaxFrameW int
	topBoxMaxFrameH int
	topBoxLayout    fbox.ItemLayout
	taskEntryLayout fbox.ItemLayout
	projEntryLayout fbox.ItemLayout
	helpLayout      fbox.ItemLayout
}

func (m *model) calculateDimensions(w, h int) tea.Cmd {
	return func() tea.Msg {
		d := &dimensions{}
		d.windowW = w
		d.windowH = h

		box := m.createFlexbox()
		rendered := box.Render(w, h)

		d.topBoxLayout = box.LayoutsHandler.GetLayout(topBoxName)
		d.taskEntryLayout = box.LayoutsHandler.GetLayout(taskEntryName)
		d.projEntryLayout = box.LayoutsHandler.GetLayout(projEntryName)
		d.helpLayout = box.LayoutsHandler.GetLayout(helpViewName)

		d.renderedW = lipgloss.Width(rendered)
		d.renderedH = lipgloss.Height(rendered)

		tb := m.createTopBox()
		d.topBoxMaxFrameW, d.topBoxMaxFrameH = tb.GetMaxItemFrameSize()
		d.projBoxW = 15
		d.projDelegMaxW = d.projBoxW - d.topBoxMaxFrameW
		d.listsH = d.topBoxLayout.ContentHeight - d.topBoxMaxFrameH
		return dimensionsCalculatedMsg{*d}
	}
}

func (m *model) logDimensions() {
	d := m.dimensions
	slog.Debug(
		"dimensions calculated",
		"focus", m.currentFocus.toString(),
		"window_width", d.windowW,
		"window_height", d.windowH,
		"proj_box_width", d.projBoxW,
		"proj_del_max_width", d.projDelegMaxW,
		"top_box_max_frame", fmt.Sprintf("%dx%d", d.topBoxMaxFrameW, d.topBoxMaxFrameH),
		"lists_height", d.listsH,
		slog.Group("rendered", "width", d.renderedW, "height", d.renderedH),
		layoutLogGrp(d.topBoxLayout),
		layoutLogGrp(d.taskEntryLayout),
		layoutLogGrp(d.projEntryLayout),
		layoutLogGrp(d.helpLayout),
	)
}

func layoutLogGrp(l fbox.ItemLayout) slog.Attr {
	return slog.Group(
		l.Name,
		"full", fmt.Sprintf("%dx%d", l.FullWidth, l.FullHeight),
		"content", fmt.Sprintf("%dx%d", l.ContentWidth, l.ContentHeight),
		"full_frame", fmt.Sprintf("%dx%d", l.FrameWidth, l.FrameHeight),
	)
}
