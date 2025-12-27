package titlebox

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Box struct {
	Title          string
	Body           string
	BoxStyle       lipgloss.Style
	TitleStyle     lipgloss.Style
	TitleAlignment Alignment
}

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

var (
	DefaultBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(0, 1)

	DefaultTitleStyle = lipgloss.NewStyle().
				Padding(0, 1)
)

func New() Box {
	return Box{
		BoxStyle:   DefaultBoxStyle,
		TitleStyle: DefaultTitleStyle,
	}
}

func (b Box) SetBoxStyle(style lipgloss.Style) Box {
	b.BoxStyle = style
	return b
}

func (b Box) SetTitleStyle(style lipgloss.Style) Box {
	b.TitleStyle = style
	return b
}

func (b Box) SetTitle(title string) Box {
	b.Title = title
	return b
}

func (b Box) SetTitleAlignment(align Alignment) Box {
	b.TitleAlignment = align
	return b
}

func (b Box) SetBody(body string) Box {
	b.Body = body
	return b
}

func (b Box) Render(w, h int) string {
	if w <= 0 || h <= 0 {
		return ""
	}

	// Get size of top border by comparing before/after removing
	// Render top row with the title
	top := topBorder(b.TitleAlignment, w, b.BoxStyle, b.TitleStyle, b.Title)
	contentH := h
	if contentH <= 0 {
		return top
	}

	// Render the rest of the box
	bottom := b.BoxStyle.
		BorderTop(false).
		Width(w).
		Height(contentH).
		Render(b.Body)

	// Stack the pieces
	return lipgloss.JoinVertical(lipgloss.Top, top, bottom)
}

func (b Box) FrameSize() (fw, fh int) {
	return b.BoxStyle.Padding(0).GetFrameSize()
}

func topBorder(align Alignment, width int, boxStyle, titleStyle lipgloss.Style, title string) string {
	border := boxStyle.GetBorderStyle()
	fg := boxStyle.GetBorderTopForeground()
	styler := lipgloss.NewStyle().Foreground(fg).Render
	topLeft := styler(border.TopLeft)
	topRight := styler(border.TopRight)
	innerW := width
	if innerW <= 0 {
		return ""
	}

	titleStr := titleStyle.Render(title)
	titleW := lipgloss.Width(titleStr)

	// If title won't fit, just render border without title
	if titleW >= innerW {
		borderLine := styler(strings.Repeat(border.Top, innerW))
		return topLeft + borderLine + topRight
	}

	leftFill := 0
	rightFill := 0
	switch align {
	case AlignCenter:
		leftFill = (innerW - titleW) / 2
		rightFill = innerW - titleW - leftFill
	case AlignRight:
		leftFill = innerW - titleW
	default:
		rightFill = innerW - titleW
	}

	left := styler(strings.Repeat(border.Top, leftFill))
	right := styler(strings.Repeat(border.Top, rightFill))
	return topLeft + left + titleStr + right + topRight
}
