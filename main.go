package main

import (
	"fmt"
	"log"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertonstz/go-workflows/components/list"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
	"github.com/evertonstz/go-workflows/shared"
)

type sessionState uint

const (
	listView sessionState = iota
	editView
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
	// BorderStyle(lipgloss.NormalBorder()).
	// BorderForeground(lipgloss.Color("69"))
)

type termDimensions struct {
	width  int
	height int
}

type model struct {
	state          sessionState
	list           list.Model
	textArea       textarea.Model
	termDimensions termDimensions
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) focused() sessionState {
	return m.state
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case shared.CopyToClipboard:
		handleClipboardCopy(msg.Desc)
	case tea.WindowSizeMsg:
		m.termDimensions.width = msg.Width
		m.termDimensions.height = msg.Height
	case shared.SaveItem:
		handleSaveItem(m, msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == listView {
				m.state = editView
				return m, func() tea.Msg {
					return shared.ItemMsg{Title: m.list.CurentItem().Title(), 
						Desc: m.list.CurentItem().Description()}
				}
			} else {
				m.state = listView
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
			focusedModelStyle.AlignHorizontal(lipgloss.Left).Width(fistPanelWidth).Height(panelHeight).Render(fmt.Sprintf("%4s", m.list.View())))
	} else {
		s = lipgloss.JoinHorizontal(lipgloss.Top,
			modelStyle.Faint(true).AlignHorizontal(lipgloss.Left).Width(fistPanelWidth).Height(panelHeight).Render(fmt.Sprintf("%4s", m.list.View())),
			focusedModelStyle.Width(secondPanelWidth).Height(panelHeight).Render(m.textArea.View()))
	}
	return s
}

func main() {
	m := model{
		list:     list.New(),
		textArea: textarea.New(),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
