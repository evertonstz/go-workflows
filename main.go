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
	"github.com/evertonstz/go-workflows/components/list"
	"github.com/evertonstz/go-workflows/components/notification"
	"github.com/evertonstz/go-workflows/components/persist"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
	"github.com/evertonstz/go-workflows/models"
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

	keys struct {
		listKeys list.KeyMap
	}

	model struct {
		confirmationModal confirmationmodal.Model
		keys              keys
		help              help.Model
		state             sessionState
		list              list.Model
		textArea          textarea.Model
		persistPath       string
		notification      notification.Model
		termDimensions    termDimensions
		currentHelpHeight int
		panelsStyle       panelsStyle
	}
	sessionState uint
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

func (m *model) updatePanelSizes() {
	currentNotificationHeight := m.panelsStyle.notificationPanelStyle.GetHeight()
	m.currentHelpHeight = strings.Count(m.help.View(m.keys.listKeys), "\n") + 1

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
		switch {
		case key.Matches(msg, m.keys.listKeys.Delete):
			m.rebuildConfirmationModel("Are you sure you want to delete this workflow?",
				"Yes",
				"No",
				shared.DeleteCurrentItemCmd(m.list.CurrentItemIndex()),
				shared.CloseConfirmationModalCmd())
			m.changeFocus(confirmationModalView)
		case key.Matches(msg, m.keys.listKeys.Help):
			m.toggleHelpShowAll()
		case key.Matches(msg, m.keys.listKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.listKeys.Esc):
			if m.state == editView {
				m.changeFocus(listView)
			}
			r, _ := m.list.Update(msg)
			m.list = r.(list.Model)
			return m, nil
		case key.Matches(msg, m.keys.listKeys.Enter):
			if m.focused() == listView && !m.list.InputOn() {
				var c tea.Cmd
				m.changeFocus(editView)

				updatedListModel, c := m.list.Update(msg)
				m.list = updatedListModel.(list.Model)
				cmds = append(cmds, c)
				updatedTextAreaModel, c := m.textArea.Update(msg)
				m.textArea = updatedTextAreaModel.(textarea.Model)
				cmds = append(cmds, c)
				return m, tea.Batch(cmds...)
			} else if m.focused() == listView && m.list.InputOn() {
				var c tea.Cmd
				_, c = m.list.Update(msg)
				return m, tea.Batch(c)
			}
		default:
			if m.help.ShowAll {
				m.toggleHelpShowAll()
			}
		}
	}

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

	notfyModel, notfyMsg := m.notification.Update(msg)
	cmds = append(cmds, notfyMsg)
	m.notification = notfyModel

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var mainContent string
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
		m.panelsStyle.helpPanelStyle.Render(m.help.View(m.keys.listKeys)))
}

func new() model {
	return model{
		confirmationModal: confirmationmodal.NewConfirmationModal("", "", "", nil, nil),
		keys:              keys{listKeys: list.Keys},
		help:              help.New(),
		list:              list.New(),
		textArea:          textarea.New(),
		notification:      notification.New("Workflows"),
		panelsStyle: panelsStyle{
			leftPanelStyle:         leftPanelStyle,
			rightPanelStyle:        rightPanelStyle,
			helpPanelStyle:         helpPanelStyle,
			notificationPanelStyle: notificationPanelStyle,
		},
		currentHelpHeight: 0,
	}
}

func main() {
	p := tea.NewProgram(new(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
