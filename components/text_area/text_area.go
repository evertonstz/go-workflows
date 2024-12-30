package textarea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
)

type Model struct {
	TextArea    textarea.Model
	currentItem models.Item
	err         error
}

func New() Model {
	ti := textarea.New()
	ti.Placeholder = "Paste or type your command here..."
	ti.Focus()

	return Model{
		TextArea: ti,
		err:      nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
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
	return fmt.Sprintf(
		"%s\n%s",
		m.TextArea.View(),
		"(esc to return to list, ctrl+s to save)",
	)
}
