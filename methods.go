package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"

	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/components/persist"
	"github.com/evertonstz/go-workflows/models"
)

func (m model) getHelpKeys() help.KeyMap {
	if m.screenState == addNew {
		return helpkeys.AddNewKeys
	}
	return helpkeys.LisKeys
}

func (m model) isSmallWidth() bool {
	return m.termDimensions.width < 100
}

func (m *model) updatePanelSizes() {
	currentNotificationHeight := m.panelsStyle.notificationPanelStyle.GetHeight()
	m.currentHelpHeight = strings.Count(m.help.View(m.getHelpKeys()), "\n") + 1

	m.addNewScreen.SetSize(m.termDimensions.width/2, m.termDimensions.height/2-(m.currentHelpHeight+currentNotificationHeight))
	m.listScreen.SetSize(m.termDimensions.width, m.termDimensions.height-(m.currentHelpHeight+currentNotificationHeight+1), m.isSmallWidth())
}

func (m *model) toggleHelpShowAll() {
	m.help.ShowAll = !m.help.ShowAll
	m.updatePanelSizes()
}

func (m model) persistItems() tea.Cmd {
	var items []models.Item
	for _, i := range m.listScreen.GetAllItems() {
		items = append(items, models.Item{
			Title:       i.Title(),
			Desc:        i.Description(),
			Command:     i.Command(),
			DateAdded:   i.DateAdded(),
			DateUpdated: i.DateUpdated()})
	}
	data := models.Items{Items: items}

	return persist.PersistListData(m.persistPath, data)
}
