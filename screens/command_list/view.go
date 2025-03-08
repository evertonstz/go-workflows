package commandlist

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
	var rightPanel string

	switch m.currentRightPanel {
	case textArea:
		rightPanel = m.panelsStyle.rightPanelStyle.Render(m.textArea.View())
	case modal:
		rightPanel = m.panelsStyle.rightPanelStyle.Render(m.confirmationModal.View())
	default:
		rightPanel = ""
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom,
		m.panelsStyle.leftPanelStyle.Render(m.list.View()),
		m.panelsStyle.rightPanelStyle.Render(rightPanel))
}
