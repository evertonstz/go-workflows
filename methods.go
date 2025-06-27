package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"

	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/shared"
	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/evertonstz/go-workflows/shared/di/services"
	"github.com/evertonstz/go-workflows/shared/messages"
)

const (
	smallWidthThreshold = 100
)

func (m model) getNotificationTitle() string {
	if m.currentPath == "/" {
		return "Workflows"
	}
	return m.currentPath
}

func (m model) getHelpKeys() help.KeyMap {
	if m.screenState == addNew {
		return helpkeys.AddNewKeys
	}
	return helpkeys.LisKeys
}

func (m model) isSmallWidth() bool {
	return m.termDimensions.width < smallWidthThreshold
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

func (m model) persistItemsV2() tea.Cmd {
	// Get database manager from DI container
	persistence := di.GetService[*services.PersistenceService](di.PersistenceServiceKey)
	databaseManager, err := services.NewDatabaseManagerV2(persistence)
	if err != nil {
		return shared.ErrorCmd(err)
	}

	database := databaseManager.GetDatabase()
	return messages.PersistListDataV2Cmd(database)
}
