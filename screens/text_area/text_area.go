package textarea

// A simple program demonstrating the textarea component from the Bubbles
// component library.

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type Model struct {
	title    string
	textarea textarea.Model
	err      error
}

func InitialModel(title string) Model {
	ti := textarea.New()
	ti.Placeholder = "Type the command you want to remember here..."
	ti.Focus()

	return Model{
		title:    title,
		textarea: ti,
		err:      nil,
	}
}

func (m *Model) SetTitle(title string) {
	m.title = title
}

func (m *Model) Value() string {
	return m.textarea.Value()
}

func (m *Model) SetValue(value string) {
	m.textarea.SetValue(value)
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf(
		m.title + "\n\n%s\n\n%s",
		m.textarea.View(),
		"(ESC to go back)",
	) + "\n\n"
}