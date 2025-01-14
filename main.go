package main

import (
	"fmt"
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
	leftPanelStyle = lipgloss.NewStyle().
			PaddingTop(2).
			AlignHorizontal(lipgloss.Left).
			BorderStyle(lipgloss.HiddenBorder())
	rightPanelStyle = lipgloss.NewStyle().
			PaddingRight(3).
			PaddingTop(3).
			Width(15).
			Height(5).
			BorderStyle(lipgloss.HiddenBorder())
)

type (
	termDimensions struct {
		width  int
		height int
	}

	model struct {
		keys           keyMap
		help           help.Model
		state          sessionState
		list           list.Model
		textArea       textarea.Model
		persistPath    string
		termDimensions termDimensions
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
		m.textArea.SetSize(int(math.Floor(float64(msg.Width)*0.5)),
			int(math.Floor(float64(msg.Height)*0.25)))
		m.list.SetSize(int(math.Floor(float64(msg.Width)*0.5)),
			int(math.Floor(float64(msg.Height)*0.75)))
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
	fistPanelWidth := int(math.Floor(float64(m.termDimensions.width) * 0.5))
	secondPanelWidth := m.termDimensions.width - fistPanelWidth
	panelHeight := m.termDimensions.height

	helpView := lipgloss.NewStyle().Width(m.termDimensions.width).PaddingLeft(2).Render(m.help.View(keys))

	helpHeight := strings.Count(helpView, "\n")

	var s string
	if m.focused() == listView {
		s = lipgloss.JoinHorizontal(lipgloss.Top,
			leftPanelStyle.
				PaddingTop(2).
				AlignHorizontal(lipgloss.Left).
				Width(fistPanelWidth).
				Height(panelHeight-helpHeight).
				Render(fmt.Sprintf("%4s", m.list.View())),
			rightPanelStyle.
				// Border(lipgloss.ThickBorder()).
				Width(secondPanelWidth).
				Height(panelHeight).
				Render(m.textArea.View()))
	} else {
		s = lipgloss.JoinHorizontal(lipgloss.Bottom,
			leftPanelStyle.
				Faint(true).
				Width(fistPanelWidth).
				Height(panelHeight).
				Render(fmt.Sprintf("%4s", m.list.View())),
			rightPanelStyle.
				// Border(lipgloss.ThickBorder()).
				Width(secondPanelWidth).
				Height(panelHeight).
				Render(m.textArea.View()))
	}
	return lipgloss.JoinVertical(lipgloss.Left, s, helpView)
}

func main() {
	m := model{
		keys:     keys,
		help:     help.New(),
		list:     list.New(),
		textArea: textarea.New(),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
