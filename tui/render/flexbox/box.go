package flexbox

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dsrosen6/yata/tui/render/titlebox"
)

type Box struct {
	Direction Direction
	Ratio     int
	Items     []Item
}

type Direction int

const (
	Vertical Direction = iota
	Horizontal
)

func New(dir Direction, ratio int) *Box {
	return &Box{
		Direction: dir,
		Ratio:     ratio,
		Items:     []Item{},
	}
}

func (b *Box) AddTitleBox(box titlebox.Box, ratio int, showFunc func() bool) *Box {
	it := TitleBoxToItem(box, ratio)
	return b.AddItem(it, showFunc)
}

func (b *Box) AddStyleBox(style lipgloss.Style, body string, ratio int, showFunc func() bool) *Box {
	it := StyleToItem(style, body, ratio)
	return b.AddItem(it, showFunc)
}

func (b *Box) AddFlexBox(box *Box, ratio int, showFunc func() bool) *Box {
	it := FlexBoxToItem(box, ratio)
	return b.AddItem(it, showFunc)
}

func (b *Box) AddItem(it Item, showFunc func() bool) *Box {
	if showFunc != nil {
		if !showFunc() {
			return b
		}
	}

	b.Items = append(b.Items, it)
	return b
}

func (b *Box) RemoveItemAt(i int) *Box {
	if i < 0 || i >= len(b.Items) {
		return b
	}

	b.Items = append(b.Items[:i], b.Items[i+1:]...)
	return b
}

func (b *Box) FrameSize() (int, int) {
	return 0, 0
}

func (b *Box) Render(w, h int) string {
	if len(b.Items) == 0 {
		return ""
	}

	// Calculate total ratio and total frame size in main axis direction
	totalRatio := 0
	totalFrameMain := 0

	for _, it := range b.Items {
		totalRatio += it.Ratio
		fw, fh := it.Node.FrameSize()
		if b.Direction == Vertical {
			totalFrameMain += fh // sum border heights
		} else {
			totalFrameMain += fw // some border widths
		}
	}

	if totalRatio <= 0 {
		return ""
	}

	mainSize := h
	crossSize := w
	if b.Direction == Horizontal {
		mainSize = w
		crossSize = h
	}

	usableMain := mainSize - totalFrameMain
	if usableMain <= 0 {
		return ""
	}

	out := make([]string, 0, len(b.Items))
	used := 0
	for i, it := range b.Items {
		fw, fh := it.Node.FrameSize()

		itemCross := crossSize
		if b.Direction == Vertical {
			itemCross -= fw
		} else {
			itemCross -= fh
		}
		if itemCross <= 0 {
			continue
		}

		itemMain := usableMain * it.Ratio / totalRatio
		if i == len(b.Items)-1 {
			itemMain = usableMain - used
		}

		if itemMain < 0 {
			itemMain = 0
		}

		used += itemMain

		if b.Direction == Vertical {
			out = append(out, it.Node.Render(itemCross, itemMain))
		} else {
			out = append(out, it.Node.Render(itemMain, itemCross))
		}
	}

	if b.Direction == Vertical {
		return lipgloss.JoinVertical(lipgloss.Top, out...)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, out...)
}
