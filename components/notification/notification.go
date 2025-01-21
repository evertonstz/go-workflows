package notification

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Msg representa uma mensagem para exibir uma notificação
type Msg struct {
	Text string
}

// Model representa o estado da notificação
type Model struct {
	Text      string
	visible   bool
	timerDone chan struct{}
}

// New cria um novo Model de notificação
func New() Model {
	return Model{
		Text:      "",
		visible:   false,
		timerDone: nil,
	}
}

// Init inicializa o modelo da notificação
func (m Model) Init() tea.Cmd {
	return nil
}

// Update atualiza o estado do modelo
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Msg:
		// Exibe a notificação e inicia o temporizador
		m.Text = msg.Text
		m.visible = true
		m.timerDone = make(chan struct{})
		return m, startTimer(2*time.Second, m.timerDone)

	case timerMsg:
		// Remove a notificação ao término do temporizador
		if m.timerDone != nil {
			close(m.timerDone)
		}
		m.Text = ""
		m.visible = false
		m.timerDone = nil
	}

	return m, nil
}

// View exibe a notificação, se visível
func (m Model) View() string {
	if !m.visible {
		return ""
	}
	return m.Text
}

// timerMsg é uma mensagem interna usada para notificar que o temporizador terminou
type timerMsg struct{}

// startTimer cria um Cmd que envia uma mensagem após um determinado período
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

// CmdShowNotification retorna um Cmd que envia uma mensagem para exibir uma notificação
func CmdShowNotification(text string) tea.Cmd {
	return func() tea.Msg {
		return Msg{Text: text}
	}
}
