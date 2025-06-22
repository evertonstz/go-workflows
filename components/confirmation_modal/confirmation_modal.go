package confirmationmodal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	localizer    *i18n.Localizer
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

func (m *Model) SetMessage(messageID string, templateData ...map[string]interface{}) {
	data := make(map[string]interface{})
	if len(templateData) > 0 {
		data = templateData[0]
	}
	m.Message = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
}

func (m *Model) SetConfirmButtonLabel(labelID string) {
	m.ConfirmButton = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: labelID})
}

func (m *Model) SetCancelButtonLabel(labelID string) {
	m.CancelButton = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: labelID})
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

func NewConfirmationModal(messageID, confirmButtonID, cancelButtonID string, confirmCmd, cancelCmd tea.Cmd, loc *i18n.Localizer) Model {
	localizer = loc
	return Model{
		Message:       localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID}),
		ConfirmButton: localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: confirmButtonID}),
		CancelButton:  localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: cancelButtonID}),
		ConfirmCmd:    confirmCmd,
		CancelCmd:     cancelCmd,
		selectedInput: confirm,
	}
}
