package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/components/notification"
	"github.com/evertonstz/go-workflows/components/persist"
	commandlist "github.com/evertonstz/go-workflows/screens/command_list"
	"github.com/evertonstz/go-workflows/shared"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case shared.ErrorMsg:
		return m, notification.ShowNotificationCmd(msg.Err.Error(), true, nil)
	case shared.DidCloseAddNewScreenMsg:
		m.screenState = newList
	case shared.DidAddNewItemMsg:
		m.screenState = newList
		updatedListModel, _ := m.listScreen.Update(msg)
		m.listScreen = updatedListModel.(commandlist.Model)
		return m, m.persistItems()
	case shared.DidDeleteItemMsg:
		cmds = append(cmds, m.persistItems())
	case shared.CopiedToClipboardMsg:
		// TODO: Add "notification_copied_to_clipboard" to en.json
		return m, notification.ShowNotificationCmd("notification_copied_to_clipboard", false, nil)
	case persist.PersistedFileMsg:
		// TODO: Add "notification_saved" to en.json
		return m, notification.ShowNotificationCmd("notification_saved", false, nil)
	case persist.InitiatedPersistion:
		m.persistPath = msg.DataFile
		return m, persist.LoadDataFileCmd(msg.DataFile)
	case tea.WindowSizeMsg:
		m.termDimensions.width = msg.Width
		m.termDimensions.height = msg.Height
		m.updatePanelSizes()
	case shared.DidUpdateItemMsg:
		return m, m.persistItems()

	case tea.KeyMsg:
		listHelpKeys := helpkeys.FullHelpKeys(localizer)
		addNewHelpKeys := helpkeys.AddNewKeys(localizer)

		switch m.screenState {
		case addNew:
			switch {
			case key.Matches(msg, addNewHelpKeys.Help):
				m.toggleHelpShowAll()
			case key.Matches(msg, addNewHelpKeys.Close): // Assuming close is used to quit from add new screen
				m.screenState = newList // or handle as quit if that's the intent
				return m, nil
			default:
				if m.help.ShowAll {
					m.toggleHelpShowAll()
				}
			}
		case newList:
			switch {
			case key.Matches(msg, listHelpKeys.AddNewWorkflow):
				m.screenState = addNew
				return m, nil
			case key.Matches(msg, listHelpKeys.Help):
				m.toggleHelpShowAll()
			case key.Matches(msg, listHelpKeys.Quit):
				return m, tea.Quit
			// TODO add edition function
			// case key.Matches(msg, listHelpKeys.Enter): // If 'Enter' is part of ListKeyMap for an action
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
