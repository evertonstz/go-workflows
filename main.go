package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	confirmationmodal "github.com/evertonstz/go-workflows/components/confirmation_modal"
	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	// "github.com/evertonstz/go-workflows/components/list"
	"github.com/evertonstz/go-workflows/components/notification"
	"github.com/evertonstz/go-workflows/components/persist"
	// textarea "github.com/evertonstz/go-workflows/components/text_area"
	"github.com/evertonstz/go-workflows/models"
	addnew "github.com/evertonstz/go-workflows/screens/add_new"
	commandlist "github.com/evertonstz/go-workflows/screens/command_list"
	"github.com/evertonstz/go-workflows/shared"
)

var (
	helpPanelStyle         = lipgloss.NewStyle().PaddingLeft(2)
	notificationPanelStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Height(1).
				AlignHorizontal(lipgloss.Left)
)

type (
	termDimensions struct {
		width  int
		height int
	}

	panelsStyle struct {
		helpPanelStyle         lipgloss.Style
		notificationPanelStyle lipgloss.Style
	}

	model struct {
		confirmationModal confirmationmodal.Model
		help              help.Model
		screenState       screenState
		addNewScreen      addnew.Model
		listScreen        commandlist.Model
		persistPath       string
		notification      notification.Model
		termDimensions    termDimensions
		currentHelpHeight int
		panelsStyle       panelsStyle
	}
	screenState  uint
)

const (
	addNew screenState = iota
	newList
)

func (m model) Init() tea.Cmd {
	return persist.InitPersistionManagerCmd("go-workflows")
}

func (m model) getHelpKeys() help.KeyMap {
	if m.screenState == addNew {
		return helpkeys.AddNewKeys
	}
	return helpkeys.LisKeys
}

func (m *model) updatePanelSizes() {
	currentNotificationHeight := m.panelsStyle.notificationPanelStyle.GetHeight()
	m.currentHelpHeight = strings.Count(m.help.View(m.getHelpKeys()), "\n") + 1

	m.addNewScreen.SetSize(m.termDimensions.width/2, m.termDimensions.height/2 - (m.currentHelpHeight + currentNotificationHeight))
	m.listScreen.SetSize(m.termDimensions.width, m.termDimensions.height - (m.currentHelpHeight + currentNotificationHeight))
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case shared.ErrorMsg:
		return m, notification.ShowNotificationCmd(msg.Err.Error())
	case shared.DidCloseAddNewScreenMsg:
		m.screenState = newList
		// return m, nil
	case shared.DidAddNewItemMsg:
		m.screenState = newList
		updatedListModel, _ := m.listScreen.Update(msg)
		m.listScreen = updatedListModel.(commandlist.Model)
		return m, m.persistItems()
	case shared.DidDeleteItemMsg:
		cmds = append(cmds, m.persistItems())
	case shared.CopiedToClipboardMsg:
		return m, notification.ShowNotificationCmd("Copied to clipboard!")
	case persist.PersistedFileMsg:
		return m, notification.ShowNotificationCmd("Saved!")
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
		if m.screenState == addNew {
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
		}

		if m.screenState == newList {
			switch {
			case key.Matches(msg, helpkeys.LisKeys.AddNewWorkflow):
				m.screenState = addNew
				return m, nil
			// case key.Matches(msg, helpkeys.LisKeys.Delete):
			// 	m.rebuildConfirmationModel("Are you sure you want to delete this workflow?",
			// 		"Yes",
			// 		"No",
			// 		shared.DeleteCurrentItemCmd(m.list.CurrentItemIndex()),
			// 		shared.CloseConfirmationModalCmd())
			// 	m.changeFocus(confirmationModalView)
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

	if m.screenState == addNew {
		addNewScreenModel, addNewScreenCmd := m.addNewScreen.Update(msg)
		cmds = append(cmds, addNewScreenCmd)
		m.addNewScreen = addNewScreenModel
	}

	if m.screenState == newList {
		var cmd tea.Cmd
		updatedListScreenModel, cmd := m.listScreen.Update(msg)
		m.listScreen = updatedListScreenModel.(commandlist.Model)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	notificationView := m.panelsStyle.notificationPanelStyle.Render(m.notification.View())
	helpView := m.panelsStyle.helpPanelStyle.Render(m.help.View(m.getHelpKeys()))

	switch m.screenState {
	case addNew:
		return lipgloss.JoinVertical(lipgloss.Left,
			notificationView,
			lipgloss.Place(m.termDimensions.width,
				m.termDimensions.height - ( m.panelsStyle.notificationPanelStyle.GetHeight() + m.currentHelpHeight),
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
}}

func new() model {
	return model{
		confirmationModal: confirmationmodal.NewConfirmationModal("", "", "", nil, nil),
		help:              help.New(),
		addNewScreen:      addnew.New(),
		listScreen:        commandlist.New(),
		notification:      notification.New("Workflows"),
		panelsStyle:       panelsStyle{
								helpPanelStyle:         helpPanelStyle,
								notificationPanelStyle: notificationPanelStyle,
							},
		currentHelpHeight:  0,
		screenState:        newList,
	}
}

func main() {
	p := tea.NewProgram(new(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
