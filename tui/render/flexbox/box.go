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

func (b *Box) AddTitleBox(box titlebox.Box, ratio int, fixedW, fixedH *int, showFunc func() bool) *Box {
	it := TitleBoxToItem(box, ratio, fixedW, fixedH)
	return b.AddItem(it, showFunc)
}

func (b *Box) AddStyleBox(s lipgloss.Style, body string, ratio int, fixedW, fixedH *int, showFunc func() bool) *Box {
	it := StyleToItem(s, body, ratio, fixedW, fixedH)
	return b.AddItem(it, showFunc)
}

func (b *Box) AddFlexBox(box *Box, ratio int, fixedW, fixedH *int, showFunc func() bool) *Box {
	it := FlexBoxToItem(box, ratio, fixedW, fixedH)
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

	// Determine main and cross sizes
	mainSize := h
	crossSize := w
	if b.Direction == Horizontal {
		mainSize = w
		crossSize = h
	}

	// First pass - calculate space allocation
	totalRatio := 0
	totalFrameMain := 0
	fixedMainTotal := 0
	lastFlexibleIdx := -1

	for i, it := range b.Items {
		fw, fh := it.Node.FrameSize()

		// add frame size for all items

		if b.Direction == Vertical {
			totalFrameMain += fh
		} else {
			totalFrameMain += fw
		}

		// check if this item has a fixed size in main direction
		hasFixedM := false
		if b.Direction == Vertical && it.FixedHeight != nil {
			fixedMainTotal += *it.FixedHeight
			hasFixedM = true
		} else if b.Direction == Horizontal && it.FixedWidth != nil {
			fixedMainTotal += *it.FixedWidth
			hasFixedM = true
		}

		if !hasFixedM {
			totalRatio += it.Ratio
			lastFlexibleIdx = i
		}
	}

	// calculate usable space for flexible items
	usableM := mainSize - totalFrameMain - fixedMainTotal
	if usableM < 0 {
		usableM = 0
	}

	// second pass: render items
	out := make([]string, 0, len(b.Items))
	usedFlexible := 0
	for i, it := range b.Items {
		fw, fh := it.Node.FrameSize()

		// calc cross size
		// calc is short for calculate
		itemCross := crossSize
		if b.Direction == Vertical {
			// cross axis is width
			if it.FixedWidth != nil {
				itemCross = *it.FixedWidth
			} else {
				itemCross -= fw
			}
		} else {
			// cross axis is height
			if it.FixedHeight != nil {
				itemCross = *it.FixedHeight
			} else {
				itemCross -= fh
			}
		}

		if itemCross <= 0 {
			continue
		}

		// calculate main size
		var itemMain int
		isFixed := false
		if b.Direction == Vertical && it.FixedHeight != nil {
			itemMain = *it.FixedHeight
			isFixed = true
		} else if b.Direction == Horizontal && it.FixedWidth != nil {
			itemMain = *it.FixedWidth
			isFixed = true
		}

		if !isFixed {
			// flexible item, use ratio
			if totalRatio > 0 {
				if i == lastFlexibleIdx {
					// last flexible index gets remaining space
					itemMain = usableM - usedFlexible
				} else {
					itemMain = usableM * it.Ratio / totalRatio
				}
				usedFlexible += itemMain
			} else {
				itemMain = 0
			}
		}

		if itemMain < 0 {
			itemMain = 0
		}

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

// FixedSize is a helper function for passing fixed sizes to constructor funcs so
// an extra line defining the int isn't necessary.
func FixedSize(s int) *int {
	return &s
}
