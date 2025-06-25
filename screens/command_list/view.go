package commandlist

import "github.com/charmbracelet/lipgloss"

func (m Model) bigWidthView() string {
	var rightPanel string

	switch m.currentRightPanel {
	case textArea:
		rightPanel = m.panelsStyle.rightPanelStyle.Render(m.textArea.View())
	case modal:
		rightPanel = m.panelsStyle.rightPanelStyle.Render(m.confirmationModal.View())
	default:
		rightPanel = ""
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.panelsStyle.leftPanelStyle.Render(m.navigableList.View()),
		m.panelsStyle.rightPanelStyle.Render(rightPanel))
}

func (m Model) smallWidthView() string {
	var rightPanel string

	switch m.currentRightPanel {
	case textArea:
		rightPanel = m.panelsStyle.rightPanelStyle.Render(m.textArea.View())
	case modal:
		rightPanel = m.panelsStyle.rightPanelStyle.Render(m.confirmationModal.View())
	default:
		rightPanel = ""
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		m.panelsStyle.leftPanelStyle.Render(m.navigableList.View()),
		m.panelsStyle.rightPanelStyle.Render(rightPanel))
}

func (m Model) View() string {
	if m.isSmallWidth {
		return m.smallWidthView()
	}
	return m.bigWidthView()
}
