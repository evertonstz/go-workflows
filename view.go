package main

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	notificationView := m.panelsStyle.notificationPanelStyle.Render(m.notification.View())
	helpView := m.panelsStyle.helpPanelStyle.Render(m.help.View(m.getHelpKeys()))

	switch m.screenState {
	case addNew:
		return lipgloss.JoinVertical(lipgloss.Left,
			notificationView,
			lipgloss.Place(m.termDimensions.width,
				m.termDimensions.height-(m.panelsStyle.notificationPanelStyle.GetHeight()+m.currentHelpHeight),
				lipgloss.Center,
				lipgloss.Center,
				m.addNewScreen.View()),
			helpView)
	case newList:
		return lipgloss.JoinVertical(lipgloss.Left,
			notificationView,
			m.listScreen.View(),
			helpView)
	default:
		return ""
	}
}
