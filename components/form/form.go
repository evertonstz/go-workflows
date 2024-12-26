package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/shared"
)

type Model struct {
	input        string
	selectedItem string
}

func New() Model {
	return Model{input: ""}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case shared.ItemMsg:
		m.selectedItem = msg.Command
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m Model) View() string {
	return "Form Input: " + m.input + "\nSelected Item: " + m.selectedItem
}
