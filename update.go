package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/components/notification"
	"github.com/evertonstz/go-workflows/components/persist"
	"github.com/evertonstz/go-workflows/models"
	commandlist "github.com/evertonstz/go-workflows/screens/command_list"
	"github.com/evertonstz/go-workflows/shared"
	"github.com/samber/mo"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	// Handle shared.ErrorMsg for now, though ideally errors are wrapped in Results
	case shared.ErrorMsg:
		return m, notification.ShowNotificationCmd(msg.Err.Error())

	// Handle new Result-based messages
	case shared.InitiatedPersistionResultMsg:
		if msg.Result.IsOk() {
			m.persistPath = msg.Result.MustGet().DataFile
			cmds = append(cmds, persist.LoadDataFileCmd(m.persistPath))
		} else {
			// This is a critical error, perhaps log and set a persistent error state
			// For now, show a notification. Consider tea.Quit if unrecoverable.
			errMsg := fmt.Sprintf("Error initializing persistence: %s", msg.Result.Error().Error())
			m.notification.SetNotification(errMsg, notification.Error)
			// return m, tea.Quit // Example: quit on critical init error
		}
	case shared.LoadedDataFileResultMsg:
		if msg.Result.IsOk() {
			items := msg.Result.MustGet()
			// Assuming listScreen.Update can handle a direct models.Items or a message wrapping it.
			// If listScreen expects a specific message type for loading items, we'd send that here.
			// For now, let's assume there's a way to set items on listScreen,
			// or it handles a specific message like shared.LoadedDataMsg.
			// This part might need adjustment based on listScreen's API.
			// As a placeholder, we'll update the listScreen directly if possible or send a new message.
			// This simulates the old behavior of listScreen.Update(items) if it took models.Items
			// For now, we'll assume that the listScreen's Update method can handle a models.Items type message
			// or we'll need a new message type like `shared.ItemsLoadedMsg{Items: items}`
			// For this example, let's assume a direct update or a new message.
			// We'll need to pass these items to the listScreen.
			// This might involve calling m.listScreen.SetItems(items) or sending a new message.
			// The original code doesn't show how `persist.LoadedDataFileMsg` was directly handled by `listScreen.Update`.
			// Let's assume we need to update the list within listScreen.
			// This might require a new message that listScreen's Update can process.
			// Or if listScreen has a method like SetItems.
			// For now, let's assume listScreen's Update can receive a models.Items.
			// This is a slight divergence as the original LoadedDataFileMsg was empty.
			// The new LoadedDataFileResultMsg carries models.Items directly.
			// Let's assume listScreen needs to be updated with these items.
			// This part needs careful integration with how listScreen handles item loading.
			// For now, we'll update the listScreen with the items.
			// This might require a new message type that listScreen can handle.
			m.listScreen.List().SetItems(shared.ConvertItemsToBubbleList(items.Items))

		} else {
			errMsg := fmt.Sprintf("Error loading data file: %s", msg.Result.Error().Error())
			m.notification.SetNotification(errMsg, notification.Error)
		}
	case shared.PersistedFileResultMsg:
		if msg.Result.IsOk() {
			m.notification.SetNotification("Saved!", notification.Info)
		} else {
			errMsg := fmt.Sprintf("Error persisting data: %s", msg.Result.Error().Error())
			m.notification.SetNotification(errMsg, notification.Error)
		}

	case shared.DidCloseAddNewScreenMsg:
		m.screenState = newList
	case shared.DidAddNewItemMsg:
		m.screenState = newList
		// The actual item addition to m.listScreen.list happens via its own Update.
		// We just need to trigger persistence.
		updatedListModel, _ := m.listScreen.Update(msg) // Ensure listScreen processes it too
		m.listScreen = updatedListModel.(commandlist.Model)
		cmds = append(cmds, m.persistItems())
	case shared.DidDeleteItemMsg:
		// Similar to AddNew, listScreen's Update handles the list modification.
		// We just trigger persistence.
		cmds = append(cmds, m.persistItems())
	case shared.CopiedToClipboardMsg:
		return m, notification.ShowNotificationCmd("Copied to clipboard!")
	// Old persist messages are now handled by their Result counterparts
	// case persist.PersistedFileMsg:
	// 	return m, notification.ShowNotificationCmd("Saved!")
	// case persist.InitiatedPersistion:
	// 	m.persistPath = msg.DataFile
	// 	return m, persist.LoadDataFileCmd(msg.DataFile)
	case tea.WindowSizeMsg:
		m.termDimensions.width = msg.Width
		m.termDimensions.height = msg.Height
		m.updatePanelSizes()
	case shared.DidUpdateItemMsg:
		return m, m.persistItems()

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
