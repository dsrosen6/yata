package flexbox

import (
	"strconv"
)

type LayoutsHandler struct {
	Layouts        []ItemLayout
	Lookup         map[string]ItemLayout
	unknownCounter int
}

func NewLayoutsHandler() *LayoutsHandler {
	return &LayoutsHandler{
		Layouts: []ItemLayout{},
		Lookup:  map[string]ItemLayout{},
	}
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

func (h *LayoutsHandler) AddLayout(l ItemLayout) {
	h.Layouts = append(h.Layouts, l)
	h.addToLookup(l)
}

func (h *LayoutsHandler) GetLayout(name string) ItemLayout {
	if l, ok := h.Lookup[name]; ok {
		return l
	}

	return ItemLayout{}
}

func (h *LayoutsHandler) addToLookup(l ItemLayout) {
	name := l.Name
	if name == "" {
		name = "unknown"
		if h.unknownCounter > 0 {
			name += strconv.Itoa(h.unknownCounter)
			h.unknownCounter++
		}
	}

	h.Lookup[name] = l
}
