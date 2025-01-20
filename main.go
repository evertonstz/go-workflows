package main

import (
	"log"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertonstz/go-workflows/components/list"
	"github.com/evertonstz/go-workflows/components/persist"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
)

var (
	leftPanelWidthPercentage = 0.5
	leftPanelStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Left).
			BorderStyle(lipgloss.NormalBorder())
	rightPanelStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Left).
			PaddingTop(0).
			Width(15).
			Height(5).
			BorderStyle(lipgloss.NormalBorder())
)

type (
	termDimensions struct {
		width  int
		height int
	}

	panelsStyle struct {
		leftPanelStyle  lipgloss.Style
		rightPanelStyle lipgloss.Style
		helpPanelStyle  lipgloss.Style
	}

	model struct {
		keys           keyMap
		help           help.Model
		state          sessionState
		list           list.Model
		textArea       textarea.Model
		persistPath    string
		termDimensions termDimensions
		helpHeight	   int
		panelsStyle    panelsStyle
	}
	sessionState uint
)

const (
	listView sessionState = iota
	editView
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
	}
	return m.state
}

func (m *model) setSizes() {
	m.helpHeight = strings.Count(m.help.View(m.keys), "\n") + 1

	m.panelsStyle.leftPanelStyle = m.panelsStyle.leftPanelStyle.
		Width(int(math.Floor(float64(m.termDimensions.width) * leftPanelWidthPercentage))).
		Height(m.termDimensions.height - m.helpHeight - 2).
		Border(lipgloss.NormalBorder())
	m.panelsStyle.rightPanelStyle = m.panelsStyle.rightPanelStyle.
		Width(m.termDimensions.width - m.panelsStyle.leftPanelStyle.GetWidth() - 4).
		Height(m.termDimensions.height - m.helpHeight - 2).
		Border(lipgloss.NormalBorder())
	m.panelsStyle.helpPanelStyle = m.panelsStyle.helpPanelStyle.Width(m.termDimensions.width).Height(m.helpHeight)

	leftWidthFrameSize, leftHeightFrameSize := m.panelsStyle.leftPanelStyle.GetFrameSize()
	rightWidthFrameSize, rightHeightFrameSize := m.panelsStyle.rightPanelStyle.GetFrameSize()

	leftPanelWidth := m.panelsStyle.leftPanelStyle.GetWidth() - leftWidthFrameSize
	rightPanelWidth := m.panelsStyle.rightPanelStyle.GetWidth() - rightWidthFrameSize

	m.list.SetSize(leftPanelWidth, m.panelsStyle.leftPanelStyle.GetHeight()-leftHeightFrameSize)
	m.textArea.SetSize(rightPanelWidth, m.panelsStyle.rightPanelStyle.GetHeight()-rightHeightFrameSize)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case persist.InitiatedPersistion:
		m.persistPath = msg.DataFile
		return m, persist.LoadDataFileCmd(msg.DataFile)
	case persist.LoadedDataFileMsg:
		m.list.Update(msg)
	case tea.WindowSizeMsg:
		m.termDimensions.width = msg.Width
		m.termDimensions.height = msg.Height
		m.setSizes()
	case shared.DidUpdateItemMsg:
		m.list.Update(msg)

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

		return m, persist.PersistListData(m.persistPath, data)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			m.setSizes()
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Esc):
			if m.state == editView {
				m.changeFocus(listView)
			}
			r, _ := m.list.Update(msg)
			m.list = r.(list.Model)
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			if m.focused() == listView && !m.list.InputOn() {
				var c tea.Cmd
				m.changeFocus(editView)

				updatedListModel, c := m.list.Update(msg)
				m.list = updatedListModel.(list.Model)
				cmds = append(cmds, c)
				updatedTextAreaModel, c := m.textArea.Update(msg)
				m.textArea = updatedTextAreaModel.(textarea.Model)
				cmds = append(cmds, c)

			} else if m.focused() == listView && m.list.InputOn() {
				_, c := m.list.Update(msg)
				cmds = append(cmds, c)
			}
		default:
			if m.help.ShowAll {
				m.help.ShowAll = false
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
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string
	if m.focused() == listView {
		s = lipgloss.JoinHorizontal(lipgloss.Bottom,
			m.panelsStyle.leftPanelStyle.Render(m.list.View()),
			m.panelsStyle.rightPanelStyle.Render(m.textArea.View()))
	} else {
		s = lipgloss.JoinHorizontal(lipgloss.Bottom,
			m.panelsStyle.leftPanelStyle.Faint(true).Render(m.list.View()),
			m.panelsStyle.rightPanelStyle.Render(m.textArea.View()))
	}
	return lipgloss.JoinVertical(lipgloss.Left, s, m.help.View(m.keys))
}

func main() {
	m := model{
		keys:     keys,
		help:     help.New(),
		list:     list.New(),
		textArea: textarea.New(),
		panelsStyle: panelsStyle{
			leftPanelStyle:  leftPanelStyle,
			rightPanelStyle: rightPanelStyle,
			helpPanelStyle:  lipgloss.NewStyle().PaddingLeft(2),
		},
		helpHeight: 0,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
