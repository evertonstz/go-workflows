package confirmationmodal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type (
	inputs uint

	Model struct {
		Message       string
		ConfirmButton string
		CancelButton  string
		ConfirmCmd    tea.Cmd
		CancelCmd     tea.Cmd
		selectedInput inputs
	}
)

const (
	confirm inputs = iota
	cancel
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SetMessage(message string) {
	m.Message = message
}

func (m *Model) SetConfirmButtonLabel(confirmButton string) {
	m.ConfirmButton = confirmButton
}

func (m *Model) SetCancelButtonLabel(cancelButton string) {
	m.CancelButton = cancelButton
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if m.selectedInput > confirm {
				m.selectedInput--
			}
		case "right", "l":
			if m.selectedInput < cancel {
				m.selectedInput++
			}
		case "enter":
			switch m.selectedInput {
			case confirm:
				return m, m.ConfirmCmd
			case cancel:
				return m, m.CancelCmd
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var confirmButton, cancelButton string

	if m.selectedInput == confirm {
		confirmButton = focusedStyle.Render("[ " + m.ConfirmButton + " ]")
		cancelButton = blurredStyle.Render("[ " + m.CancelButton + " ]")
	} else {
		confirmButton = blurredStyle.Render("[ " + m.ConfirmButton + " ]")
		cancelButton = focusedStyle.Render("[ " + m.CancelButton + " ]")
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.Message,
		lipgloss.JoinHorizontal(lipgloss.Center, confirmButton, " ", cancelButton),
	)
}

func NewConfirmationModal(message, confirmButton, cancelButton string, confirmCmd, cancelCmd tea.Cmd) Model {
	return Model{
		Message:       message,
		ConfirmButton: confirmButton,
		CancelButton:  cancelButton,
		ConfirmCmd:    confirmCmd,
		CancelCmd:     cancelCmd,
		selectedInput: confirm,
	}
}