package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/components/notification"
	commandlist "github.com/evertonstz/go-workflows/screens/command_list"
	"github.com/evertonstz/go-workflows/shared"
	"github.com/evertonstz/go-workflows/shared/messages"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case shared.ErrorMsg:
		return m, notification.ShowNotificationCmd(msg.Err.Error())
	case shared.DidCloseAddNewScreenMsg:
		m.screenState = newList
	case shared.DidAddNewItemMsg:
		m.screenState = newList
		updatedListModel, _ := m.listScreen.Update(msg)
		m.listScreen = updatedListModel.(commandlist.Model)
		return m, m.persistItemsV2()
	case shared.DidDeleteItemMsg:
		updatedListModel, _ := m.listScreen.Update(msg)
		m.listScreen = updatedListModel.(commandlist.Model)
		return m, m.persistItemsV2()
	case shared.DidNavigateToFolderMsg:
		m.currentPath = msg.Path
		m.notification.SetDefaultText(m.getNotificationTitle())
		updatedListModel, _ := m.listScreen.Update(msg)
		m.listScreen = updatedListModel.(commandlist.Model)
		return m, nil
	case shared.CopiedToClipboardMsg:
		return m, notification.ShowNotificationCmd("Copied to clipboard!")
	case messages.PersistedFileV2Msg:
		return m, notification.ShowNotificationCmd("Saved!")
	case messages.PersistedFileMsg:
		return m, notification.ShowNotificationCmd("Saved!")
	case messages.InitiatedPersistionMsg:
		m.persistPath = msg.DataFile
		return m, messages.LoadDataFileV2Cmd()
	case messages.LoadedDataFileV2Msg:
		updatedListModel, _ := m.listScreen.Update(msg)
		m.listScreen = updatedListModel.(commandlist.Model)
		return m, nil
	case tea.WindowSizeMsg:
		m.termDimensions.width = msg.Width
		m.termDimensions.height = msg.Height
		m.updatePanelSizes()
	case shared.DidUpdateItemMsg:
		updatedListModel, _ := m.listScreen.Update(msg)
		m.listScreen = updatedListModel.(commandlist.Model)
		return m, m.persistItemsV2()
	case shared.DidSetCurrentItemMsg:
		updatedListModel, _ := m.listScreen.Update(msg)
		m.listScreen = updatedListModel.(commandlist.Model)
		return m, nil
	case shared.DidSetCurrentFolderMsg:
		updatedListModel, _ := m.listScreen.Update(msg)
		m.listScreen = updatedListModel.(commandlist.Model)
		return m, nil

	case tea.KeyMsg:
		switch m.screenState {
		case addNew:
			switch {
			case key.Matches(msg, helpkeys.LisKeys.Help):
				m.toggleHelpShowAll()
			case key.Matches(msg, helpkeys.LisKeys.Quit):
				return m, tea.Quit
			default:
				if m.help.ShowAll {
					m.toggleHelpShowAll()
				}
			}
		case newList:
			switch {
			case key.Matches(msg, helpkeys.LisKeys.AddNewWorkflow):
				m.screenState = addNew
				return m, nil
			case key.Matches(msg, helpkeys.LisKeys.Help):
				m.toggleHelpShowAll()
			case key.Matches(msg, helpkeys.LisKeys.Quit):
				return m, tea.Quit
			case key.Matches(msg, helpkeys.LisKeys.Esc):
				if m.listScreen.IsAtRoot() {
					return m, tea.Quit
				}
			// TODO add edition function
			// case key.Matches(msg, m.keys.listKeys.Enter):
			// 	m.addNewScreen.SetValues(m.list.CurentItem().Title(),
			// 								m.list.CurentItem().Description(),
			// 								m.list.CurentItem().Command())
			// 	m.screen = addNew
			// 	return m, nil
			default:
				if m.help.ShowAll {
					m.toggleHelpShowAll()
				}
			}
		}
	}

	notfyModel, notfyCmd := m.notification.Update(msg)
	cmds = append(cmds, notfyCmd)
	m.notification = notfyModel

	switch m.screenState {
	case addNew:
		addNewScreenModel, addNewScreenCmd := m.addNewScreen.Update(msg)
		cmds = append(cmds, addNewScreenCmd)
		m.addNewScreen = addNewScreenModel
	case newList:
		var cmd tea.Cmd
		updatedListScreenModel, cmd := m.listScreen.Update(msg)
		m.listScreen = updatedListScreenModel.(commandlist.Model)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
