package main

import (
	"log"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	confirmationmodal "github.com/evertonstz/go-workflows/components/confirmation_modal"
	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/components/list"
	"github.com/evertonstz/go-workflows/components/notification"
	"github.com/evertonstz/go-workflows/components/persist"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
	"github.com/evertonstz/go-workflows/models"
	addnew "github.com/evertonstz/go-workflows/screens/add_new"
	"github.com/evertonstz/go-workflows/shared"
)

var (
	leftPanelWidthPercentage = 0.5
	leftPanelStyle           = lipgloss.NewStyle().
					AlignHorizontal(lipgloss.Left)
	rightPanelStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Left).
			PaddingTop(2).
			Width(15).
			Height(5)
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
		leftPanelStyle         lipgloss.Style
		rightPanelStyle        lipgloss.Style
		helpPanelStyle         lipgloss.Style
		notificationPanelStyle lipgloss.Style
	}

	model struct {
		confirmationModal confirmationmodal.Model
		help              help.Model
		state             sessionState
		screen            screenState
		list              list.Model
		textArea          textarea.Model
		addNewScreen      addnew.Model
		persistPath       string
		notification      notification.Model
		termDimensions    termDimensions
		currentHelpHeight int
		panelsStyle       panelsStyle
	}
	sessionState uint
	screenState  uint
)

const (
	addNew screenState = iota
	listScreen
)

const (
	listView sessionState = iota
	editView
	confirmationModalView
)

func (m model) Init() tea.Cmd {
	return persist.InitPersistionManagerCmd("go-workflows")
}

func (m model) focused() sessionState {
	return m.state
}

func (m *model) changeFocus(v sessionState) sessionState {
	m.state = v
	switch v {
	case listView:
		m.textArea.SetEditing(false)
	case editView:
		m.textArea.SetEditing(true)
	case confirmationModalView:
		m.textArea.SetEditing(false)
	}
	return m.state
}

func (m model) getHelpKeys() help.KeyMap {
	if m.screen == addNew {
		return helpkeys.AddNewKeys
	}
	return helpkeys.LisKeys
}

func (m *model) updatePanelSizes() {
	currentNotificationHeight := m.panelsStyle.notificationPanelStyle.GetHeight()
	m.currentHelpHeight = strings.Count(m.help.View(m.getHelpKeys()), "\n") + 1

	m.panelsStyle.leftPanelStyle = m.panelsStyle.leftPanelStyle.
		Width(int(math.Floor(float64(m.termDimensions.width) * leftPanelWidthPercentage))).
		Height(m.termDimensions.height - m.currentHelpHeight - currentNotificationHeight)
	m.panelsStyle.rightPanelStyle = m.panelsStyle.rightPanelStyle.
		Width(m.termDimensions.width - m.panelsStyle.leftPanelStyle.GetWidth() - 4).
		Height(m.termDimensions.height - m.currentHelpHeight - currentNotificationHeight)
	m.panelsStyle.helpPanelStyle = m.panelsStyle.helpPanelStyle.
		Width(m.termDimensions.width).
		Height(m.currentHelpHeight)

	leftWidthFrameSize, leftHeightFrameSize := m.panelsStyle.leftPanelStyle.GetFrameSize()
	rightWidthFrameSize, rightHeightFrameSize := m.panelsStyle.rightPanelStyle.GetFrameSize()

	leftPanelWidth := m.panelsStyle.leftPanelStyle.GetWidth() - leftWidthFrameSize
	rightPanelWidth := m.panelsStyle.rightPanelStyle.GetWidth() - rightWidthFrameSize

	m.list.SetSize(leftPanelWidth, m.panelsStyle.leftPanelStyle.GetHeight()-leftHeightFrameSize)
	m.textArea.SetSize(rightPanelWidth, m.panelsStyle.rightPanelStyle.GetHeight()-rightHeightFrameSize)
	m.addNewScreen.SetSize(m.termDimensions.width/2, m.termDimensions.height/2-leftHeightFrameSize-m.currentHelpHeight)
}

func (m *model) toggleHelpShowAll() {
	m.help.ShowAll = !m.help.ShowAll
	m.updatePanelSizes()
}

func (m *model) rebuildConfirmationModel(title string, confirm string, cancel string, confirmCmd tea.Cmd, cancelCmd tea.Cmd) {
	m.confirmationModal = confirmationmodal.NewConfirmationModal(
		title,
		confirm,
		cancel,
		confirmCmd,
		cancelCmd,
	)
}

func (m model) persistItems() tea.Cmd {
	var items []models.Item
	for _, i := range m.list.AllItems() {
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
		m.screen = listScreen
		m.changeFocus(listView)
		return m, nil
	case shared.DidAddNewItemMsg:
		m.screen = listScreen
		updatedListModel, _ := m.list.Update(msg)
		m.list = updatedListModel.(list.Model)
		return m, m.persistItems()
	case shared.DidCloseConfirmationModalMsg:
		m.changeFocus(listView)
	case shared.DidDeleteItemMsg:
		updatedListModel, cmd := m.list.Update(msg)
		m.list = updatedListModel.(list.Model)
		cmds = append(cmds, cmd)
		m.changeFocus(listView)
		cmds = append(cmds, m.persistItems())

		return m, tea.Batch(cmds...)
	case shared.CopiedToClipboardMsg:
		return m, notification.ShowNotificationCmd("Copied to clipboard!")
	case persist.PersistedFileMsg:
		return m, notification.ShowNotificationCmd("Saved!")
	case persist.InitiatedPersistion:
		m.persistPath = msg.DataFile
		return m, persist.LoadDataFileCmd(msg.DataFile)
	case persist.LoadedDataFileMsg:
		m.list.Update(msg)
	case tea.WindowSizeMsg:
		m.termDimensions.width = msg.Width
		m.termDimensions.height = msg.Height
		m.updatePanelSizes()
	case shared.DidUpdateItemMsg:
		m.list.Update(msg)
		return m, m.persistItems()

	case tea.KeyMsg:
		if m.screen == addNew {
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

		if m.screen == listScreen {
			switch {
			case key.Matches(msg, helpkeys.LisKeys.AddNewWorkflow):
				m.screen = addNew
				return m, nil
			case key.Matches(msg, helpkeys.LisKeys.Delete):
				m.rebuildConfirmationModel("Are you sure you want to delete this workflow?",
					"Yes",
					"No",
					shared.DeleteCurrentItemCmd(m.list.CurrentItemIndex()),
					shared.CloseConfirmationModalCmd())
				m.changeFocus(confirmationModalView)
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

	if m.screen == addNew {
		addNewScreenModel, addNewScreenCmd := m.addNewScreen.Update(msg)
		cmds = append(cmds, addNewScreenCmd)
		m.addNewScreen = addNewScreenModel
	}

	if m.screen == listScreen {
		switch m.focused() {
		case listView:
			var c tea.Cmd
			var updatedListAreaModel tea.Model
			var updatedTextAreaModel tea.Model
			updatedListAreaModel, c = m.list.Update(msg)
			m.list = updatedListAreaModel.(list.Model)
			cmds = append(cmds, c)
			updatedTextAreaModel, c = m.textArea.Update(msg)
			m.textArea = updatedTextAreaModel.(textarea.Model)
			cmds = append(cmds, c)
		case editView:
			var c tea.Cmd
			var updatedTextAreaModel tea.Model
			updatedTextAreaModel, c = m.textArea.Update(msg)
			m.textArea = updatedTextAreaModel.(textarea.Model)
			cmds = append(cmds, c)
		case confirmationModalView:
			var c tea.Cmd
			var updatedConfirmationModalModel tea.Model
			updatedConfirmationModalModel, c = m.confirmationModal.Update(msg)
			m.confirmationModal = updatedConfirmationModalModel.(confirmationmodal.Model)
			cmds = append(cmds, c)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var mainContent string

	if m.screen == addNew {
		return lipgloss.JoinVertical(lipgloss.Left,
			m.panelsStyle.notificationPanelStyle.Render(m.notification.View()),
			lipgloss.Place(m.termDimensions.width,
				m.panelsStyle.leftPanelStyle.GetHeight(),
				lipgloss.Center,
				lipgloss.Center,
				m.addNewScreen.View()),
			m.panelsStyle.helpPanelStyle.Render(m.help.View(m.getHelpKeys())))
	}

	if m.focused() == listView {
		mainContent = lipgloss.JoinHorizontal(lipgloss.Bottom,
			m.panelsStyle.leftPanelStyle.Render(m.list.View()),
			m.panelsStyle.rightPanelStyle.Render(m.textArea.View()))
	} else if m.focused() == editView {
		mainContent = lipgloss.JoinHorizontal(lipgloss.Bottom,
			m.panelsStyle.leftPanelStyle.Faint(true).Render(m.list.View()),
			m.panelsStyle.rightPanelStyle.Render(m.textArea.View()))
	} else if m.focused() == confirmationModalView {
		mainContent = lipgloss.JoinHorizontal(lipgloss.Bottom,
			m.panelsStyle.leftPanelStyle.Faint(true).Render(m.list.View()),
			m.panelsStyle.rightPanelStyle.Render(m.confirmationModal.View()))
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		m.panelsStyle.notificationPanelStyle.Render(m.notification.View()),
		mainContent,
		m.panelsStyle.helpPanelStyle.Render(m.help.View(m.getHelpKeys())))
}

func new() model {
	return model{
		confirmationModal: confirmationmodal.NewConfirmationModal("", "", "", nil, nil),
		help:              help.New(),
		list:              list.New(),
		textArea:          textarea.New(),
		addNewScreen:      addnew.New(),
		notification:      notification.New("Workflows"),
		panelsStyle: panelsStyle{
			leftPanelStyle:         leftPanelStyle,
			rightPanelStyle:        rightPanelStyle,
			helpPanelStyle:         helpPanelStyle,
			notificationPanelStyle: notificationPanelStyle,
		},
		currentHelpHeight: 0,
		screen:            listScreen,
	}
}

func main() {
	p := tea.NewProgram(new(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
