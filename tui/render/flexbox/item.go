package flexbox

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dsrosen6/yata/tui/render/titlebox"
)

type Item struct {
	Name        string
	Ratio       int
	FixedWidth  *int
	FixedHeight *int
	Node        Node
}

// ItemLayout shows the width and height of an item in the flexbox, which will
// sometimes be needed for calculations outside of the flexbox itself.
type ItemLayout struct {
	Name          string
	ContentWidth  int
	ContentHeight int
	FrameWidth    int
	FrameHeight   int
	FullWidth     int
	FullHeight    int
}

// Node represents a renderable component in the flexbox layout.
//
// FrameSize should return the border dimensions only (excluding padding).
// This is because lipgloss.Style.Width() treats the width parameter as
// including padding but excluding borders.
//
// Render(w, h) receives dimensions that include padding but exclude borders,
// and should output content of total size w+borderWidth by n+borderHeight.
type Node interface {
	Render(w, h int) string
	FrameSize() (int, int)
}

// StyleNode is a standard node that can be created with just a lipgloss style and body.
type StyleNode struct {
	Style lipgloss.Style
	Body  string
}

func TitleBoxToItem(box titlebox.Box, name string, ratio int, fixedW, fixedH *int) Item {
	return Item{
		Name:        name,
		Ratio:       ratio,
		FixedWidth:  fixedW,
		FixedHeight: fixedH,
		Node:        box,
	}
}

func StyleToItem(style lipgloss.Style, name, body string, ratio int, fixedW, fixedH *int) Item {
	return Item{
		Name:        name,
		Ratio:       ratio,
		FixedWidth:  fixedW,
		FixedHeight: fixedH,
		Node:        NewStyleNode(style, body),
	}
}

func FlexBoxToItem(box *Box, name string, ratio int, fixedW, fixedH *int) Item {
	return Item{
		Name:        name,
		Ratio:       ratio,
		FixedWidth:  fixedW,
		FixedHeight: fixedH,
		Node:        box,
	}
}

func NewStyleNode(style lipgloss.Style, body string) StyleNode {
	return StyleNode{
		Style: style,
		Body:  body,
	}
}

func (sn StyleNode) Render(w, h int) string {
	if w <= 0 || h <= 0 {
		return ""
	}

	return sn.Style.
		Width(w).
		Height(h).
		Render(sn.Body)
}

func (sn StyleNode) FrameSize() (int, int) {
	return sn.Style.Padding(0).GetFrameSize()
}
