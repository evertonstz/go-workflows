package notification

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Msg struct {
	Text string
}

type Model struct {
	Text        string
	defaultText string
	visible     bool
	timerDone   chan struct{}
}

func New(defaultText string) Model {
	return Model{
		Text:        "",
		defaultText: defaultText,
		visible:     false,
		timerDone:   nil,
	}
}

var style = lipgloss.NewStyle().
				PaddingLeft(1).
				PaddingRight(1).	
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("230"))

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Msg:
		m.Text = msg.Text
		m.visible = true
		m.timerDone = make(chan struct{})
		return m, startTimer(2*time.Second, m.timerDone)

	case timerMsg:
		if m.timerDone != nil {
			close(m.timerDone)
		}
		m.Text = ""
		m.visible = false
		m.timerDone = nil
	}

	return m, nil
}

func (m Model) View() string {
	if !m.visible {
		return style.SetString(m.defaultText).Render()
	}
	return style.SetString(m.Text).Render()
}

type timerMsg struct{}

func startTimer(duration time.Duration, done chan struct{}) tea.Cmd {
	return func() tea.Msg {
		select {
		case <-time.After(duration):
			return timerMsg{}
		case <-done:
			return nil
		}
	}
}

func CmdShowNotification(text string) tea.Cmd {
	return func() tea.Msg {
		return Msg{Text: text}
	}
}
