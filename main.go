package main

import (
	"fmt"
	"log"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertonstz/go-workflows/components/list"
	"github.com/evertonstz/go-workflows/components/persist"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
	"github.com/evertonstz/go-workflows/shared"
)

var (
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.HiddenBorder())
)

type (
	termDimensions struct {
		width  int
		height int
	}

	model struct {
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
	return persist.InitPersistionManager("go-workflows")
}

func (m model) focused() sessionState {
	return m.state
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case persist.Paths:
		m.persistPath = msg.DataFile
		return m, persist.LoadPersistionFile(msg.DataFile)
	case persist.PersistionFileLoadedMsg:
		m.list.Update(m)
	case shared.CopyToClipboardMsg:
		return handleClipboardCopy(m, msg)
	case tea.WindowSizeMsg:
		m.termDimensions.width = msg.Width
		m.termDimensions.height = msg.Height
	case shared.SaveCommandMsg:
		return handleSaveItem(m, msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.state == listView && !m.list.InputOn() {
				m.state = editView
				return m, func() tea.Msg {
					return shared.ItemMsg{Title: m.list.CurentItem().Title(),
						Command: m.list.CurentItem().Command()}
				}
			} else if m.state == listView && m.list.InputOn() {
				var c tea.Cmd
				_, c = m.list.Update(msg)
				return m, tea.Batch(c)
			}
		}
	}

	switch m.focused() {
	case listView:
		var c tea.Cmd
		var updatedListAreaModel tea.Model
		updatedListAreaModel, c = m.list.Update(msg)
		m.list = updatedListAreaModel.(list.Model)
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
	var s string

	if m.focused() == listView {
		s = lipgloss.JoinHorizontal(lipgloss.Top,
			focusedModelStyle.
				AlignHorizontal(lipgloss.Left).
				Width(fistPanelWidth).
				Height(panelHeight).
				Render(fmt.Sprintf("%4s", m.list.View())))
	} else {
		s = lipgloss.JoinHorizontal(lipgloss.Top,
			modelStyle.Faint(true).
				AlignHorizontal(lipgloss.Left).
				Width(fistPanelWidth).
				Height(panelHeight).
				Render(fmt.Sprintf("%4s", m.list.View())),
			focusedModelStyle.
				Width(secondPanelWidth).
				Height(panelHeight).
				Render(m.textArea.View()))
	}
	return s
}

func main() {
	m := model{
		list:     list.New(),
		textArea: textarea.New(),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
