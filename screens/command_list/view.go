package commandlist

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Bottom,
		m.panelsStyle.leftPanelStyle.Render(m.list.View()),
		m.panelsStyle.rightPanelStyle.Render(m.textArea.View()))
}
