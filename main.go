package main

import (
	"fmt"
	"log"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertonstz/go-workflows/components/footer"
	"github.com/evertonstz/go-workflows/components/list"
	"github.com/evertonstz/go-workflows/components/text_area"
)

type sessionState uint

const (
	leftView sessionState = iota
	rightView
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
	list           tea.Model
	textArea       tea.Model
	footer         footer.Model
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
	case tea.WindowSizeMsg:
		m.termDimensions.width = msg.Width
		m.termDimensions.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == leftView {
				m.state = rightView
			} else {
				m.state = leftView
			}
		}
	}

	switch m.focused() {
	case leftView:
		var c tea.Cmd
		m.list, c = m.list.Update(msg)
		cmds = append(cmds, c)
	case rightView:
		var c tea.Cmd
		m.textArea, c = m.textArea.Update(msg)
		cmds = append(cmds, c)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	fistPanelWidth := int(math.Floor(float64(m.termDimensions.width) * 0.5))
	secondPanelWidth := m.termDimensions.width - fistPanelWidth
	panelHeight := m.termDimensions.height
	var s string
	if m.focused() == leftView {
		s = lipgloss.JoinHorizontal(lipgloss.Top, 
			focusedModelStyle.AlignHorizontal(lipgloss.Left).Width(fistPanelWidth).Height(panelHeight).Render(fmt.Sprintf("%4s", m.list.View())), 
			modelStyle.Faint(true).Width(secondPanelWidth).Height(panelHeight).Render(m.textArea.View()))
	} else {
		s = lipgloss.JoinHorizontal(lipgloss.Top, 
			modelStyle.Faint(true).AlignHorizontal(lipgloss.Left).Width(fistPanelWidth).Height(panelHeight).Render(fmt.Sprintf("%4s", m.list.View())), 
			focusedModelStyle.Width(secondPanelWidth).Height(panelHeight).Render(m.textArea.View()))
	}
	return s
}

func main() {
	m := model{
		list:   list.New(),
		textArea: textarea.New(),
		footer: footer.New(),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
