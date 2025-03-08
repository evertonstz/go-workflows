package commandlist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/components/list"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	var listModel tea.Model
	listModel, cmd = m.list.Update(msg)
	m.list = listModel.(list.Model)
	cmds = append(cmds, cmd)
	taModel, cmd := m.textArea.Update(msg)
	m.textArea = taModel.(textarea.Model)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
