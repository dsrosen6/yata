package flexbox

import (
	"slices"

	"github.com/charmbracelet/lipgloss"
	"github.com/dsrosen6/yata/tui/render/titlebox"
)

type Box struct {
	Direction      Direction
	Ratio          int
	Items          []Item
	LayoutsHandler *LayoutsHandler
}

type Direction int

const (
	Vertical Direction = iota
	Horizontal
)

func New(dir Direction, ratio int) *Box {
	return &Box{
		Direction:      dir,
		Ratio:          ratio,
		Items:          []Item{}, // slice used for rendering items in correct order
		LayoutsHandler: NewLayoutsHandler(),
	}
}

func (b *Box) AddTitleBox(box titlebox.Box, name string, ratio int, fixedW, fixedH *int, showFunc func() bool) *Box {
	it := TitleBoxToItem(box, name, ratio, fixedW, fixedH)
	return b.AddItem(it, showFunc)
}

func (b *Box) AddStyleBox(s lipgloss.Style, name, body string, ratio int, fixedW, fixedH *int, showFunc func() bool) *Box {
	it := StyleToItem(s, name, body, ratio, fixedW, fixedH)
	return b.AddItem(it, showFunc)
}

func (b *Box) AddFlexBox(box *Box, name string, ratio int, fixedW, fixedH *int, showFunc func() bool) *Box {
	it := FlexBoxToItem(box, name, ratio, fixedW, fixedH)
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

// GetAllItemsFrameSize gets the combined total frame widths and heights of a flexbox.
func (b *Box) GetAllItemsFrameSize() (int, int) {
	var w, h int
	for _, it := range b.Items {
		fw, fh := it.Node.FrameSize()
		w += fw
		h += fh
	}

	return w, h
}

// GetMaxItemFrameSize gets the maximum horizontal and vertical frame sizes from all items in a flexbox.
func (b *Box) GetMaxItemFrameSize() (int, int) {
	var widths, heights []int
	for _, it := range b.Items {
		fw, fh := it.Node.FrameSize()
		widths = append(widths, fw)
		heights = append(heights, fh)
	}
	w := slices.Max(widths)
	h := slices.Max(heights)
	return w, h
}

// Render calculates all of a box's current item layouts, and then returns the rendered string.
func (b *Box) Render(w, h int) string {
	if len(b.Items) == 0 {
		return ""
	}

	b.LayoutsHandler = b.calculateItemLayouts(w, h)

	out := make([]string, 0, len(b.Items))
	for i, it := range b.Items {
		layout := b.LayoutsHandler.Layouts[i]
		if layout.ContentWidth <= 0 || layout.ContentHeight <= 0 {
			continue
		}

		out = append(out, it.Node.Render(layout.ContentWidth, layout.ContentHeight))
	}

	if b.Direction == Vertical {
		return lipgloss.JoinVertical(lipgloss.Top, out...)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, out...)
}

// GetItemLayout is a shortcut method to calculate the box's layouts and then return
// one specifically by name.
func (b *Box) GetItemLayout(name string, w, h int) ItemLayout {
	if b.LayoutsHandler == nil {
		b.LayoutsHandler = b.CalculateItemLayouts(w, h)
	}

	return b.LayoutsHandler.GetLayout(name)
}

// CalculateItemLayouts calculates all of a box's current item layouts and then returns
// the layout handler.
func (b *Box) CalculateItemLayouts(w, h int) *LayoutsHandler {
	return b.calculateItemLayouts(w, h)
}

func (b *Box) calculateItemLayouts(w, h int) *LayoutsHandler {
	lh := NewLayoutsHandler()
	if len(b.Items) == 0 {
		return lh
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
	usableM = max(0, usableM)

	// second pass: render items
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
			lh.AddLayout(ItemLayout{Name: it.Name})
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

		itemMain = max(0, itemMain)

		cw := itemMain
		ch := itemCross
		if b.Direction == Vertical {
			cw = itemCross
			ch = itemMain
		}

		l := ItemLayout{
			Name:          it.Name,
			ContentWidth:  cw,
			ContentHeight: ch,
			FrameWidth:    fw,
			FrameHeight:   fh,
			FullWidth:     cw + fw,
			FullHeight:    ch + fh,
		}

		lh.AddLayout(l)
	}

	return lh
}

// FixedSize is a helper function for passing fixed sizes to constructor funcs so
// an extra line defining the int isn't necessary.
func FixedSize(s int) *int {
	return &s
}
