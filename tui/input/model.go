package input

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	Options struct {
		FocusedColor    string
		UnfocusedColors string
		CursorMode      cursor.Mode
		FieldKeys       []string
	}

	Model struct {
		Inputs []textinput.Model
		State  State
		Result Result

		focusedStyle   lipgloss.Style
		unfocusedStyle lipgloss.Style
		cursor         cursor.Model
		focusIndex     int
	}

	Result    map[string]string
	ResultMsg struct{ result Result }
	State     int
)

const (
	StateActive State = iota
	StateDone
)

var ErrNoFields = errors.New("no input fields provided")

func InitialInputModel(c *Options) (*Model, error) {
	if len(c.FieldKeys) == 0 {
		return nil, ErrNoFields
	}

	f := append([]string{}, c.FieldKeys...)
	m := &Model{
		Inputs:         make([]textinput.Model, len(f)),
		Result:         make(Result, len(c.FieldKeys)),
		focusedStyle:   focusedStyle(c.FocusedColor),
		unfocusedStyle: unfocusedStyle(c.UnfocusedColors),
		cursor:         makeCursor(c.CursorMode, focusedStyle(c.FocusedColor)),
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.Cursor = m.cursor
		t.Cursor.Style = m.focusedStyle
		t.Prompt = fmt.Sprintf("%s > ", f[i])
		t.CharLimit = 32

		if i == 0 {
			t.Focus()
			t.PromptStyle = m.focusedStyle
			t.TextStyle = m.focusedStyle
		}

		m.Inputs[i] = t
	}

	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			// TODO: how do we want to handle this as a nested model?
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.Inputs)-1 {
				return m, inputResultCmd(m.Inputs)
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := 0; i <= len(m.Inputs)-1; i++ {
				// if focused index, give focused style
				if i == m.focusIndex {
					cmds[i] = m.Inputs[i].Focus()
					m.Inputs[i].PromptStyle = m.focusedStyle
					m.Inputs[i].TextStyle = m.focusedStyle
					continue
				}
				// otherwise, remove style
				m.Inputs[i].Blur()
				m.Inputs[i].PromptStyle = m.unfocusedStyle
				m.Inputs[i].TextStyle = m.unfocusedStyle
			}

			return m, tea.Batch(cmds...)
		}

	case ResultMsg:
		m.State = StateDone
		m.Result = msg.result
		return m, nil
	}

	// handle character input and blinking
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *Model) View() string {
	var b strings.Builder

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	// only focused inputs will respond, so it's fine to update all
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func makeCursor(mode cursor.Mode, style lipgloss.Style) cursor.Model {
	c := cursor.New()
	if mode != 0 {
		c.SetMode(mode)
	}

	c.Style = style
	return c
}

func inputResultCmd(inputs []textinput.Model) tea.Cmd {
	return func() tea.Msg {
		r := make(Result, len(inputs))
		for _, i := range inputs {
			k := i.Prompt
			r[k] = i.Value()
		}

		return ResultMsg{result: r}
	}
}
