package textarea

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
)

type Model struct {
	TextArea        textarea.Model
	highlightedText string
	currentItem     models.Item
	editing         bool
	err             error
}

func (m *Model) SetEditing(editing bool) {
	m.editing = editing
}

func New() Model {
	ti := textarea.New()
	ti.ShowLineNumbers = false
	ti.Placeholder = "Paste or type your command here..."
	ti.Focus()
	ti.Prompt = ""

	return Model{
		TextArea:        ti,
		highlightedText: ti.Placeholder,
		editing:         false,
		err:             nil,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SetSize(width, height int) {
	m.TextArea.SetWidth(width)
	m.TextArea.SetHeight(height)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case shared.DidSetCurrentItemMsg:
		m.currentItem = msg.Item
		m.TextArea.SetValue(m.currentItem.Command)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.TextArea.Focused() {
				m.TextArea.Blur()
			}
		case tea.KeyCtrlS:
			m.currentItem.Command = m.TextArea.Value()
			return m, shared.UpdateItemCmd(m.currentItem)
		default:
			if !m.TextArea.Focused() {
				cmd = m.TextArea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}

	m.TextArea, cmd = m.TextArea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.editing {
		return m.TextArea.View()
	}

	rawText := m.TextArea.Value()
	highlightedText := SyntaxHighlight(rawText)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		highlightedText,
	)
}
