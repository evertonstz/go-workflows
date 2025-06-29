package addnew

import (
	"github.com/charmbracelet/bubbles/key"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/evertonstz/go-workflows/components/notification"
	"github.com/evertonstz/go-workflows/shared"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	titleModel, titleCmd := m.Title.Update(msg)
	descModel, descCmd := m.Description.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Down):
			switch m.selectedInput {
			case title:
				return m.focusInput(description)
			case description:
				return m.focusInput(textArea)
			case textArea:
				return m.focusInput(submit)
			case submit, close:
				return m, nil
			}
		case key.Matches(msg, m.Keys.Up):
			switch m.selectedInput {
			case title:
				return m, nil
			case description:
				return m.focusInput(title)
			case textArea:
				return m.focusInput(description)
			case submit, close:
				return m.focusInput(textArea)
			}
		case key.Matches(msg, m.Keys.Right):
			if m.selectedInput == submit {
				return m.focusInput(close)
			}
		case key.Matches(msg, m.Keys.Left):
			if m.selectedInput == close {
				return m.focusInput(submit)
			}
		case key.Matches(msg, m.Keys.Close):
			m.ResetForm()
			return m, shared.CloseAddNewScreenCmd()
		case key.Matches(msg, m.Keys.Submit):
			switch m.selectedInput {
			case submit:
				if m.isFormValid() {
					var title, description, command string
					title = m.Title.Value()
					description = m.Description.Value()
					command = m.TextArea.Value()

					m.ResetForm()
					return m, shared.AddNewItemCmd(title, description, command)
				}
				return m, notification.ShowNotificationCmd(m.notifications.fillAllFields)
			case close:
				m.ResetForm()
				return m, shared.CloseAddNewScreenCmd()
			}
		}
	}
	textModel, textCmd := m.TextArea.Update(msg)
	return Model{
		Title:         titleModel,
		Description:   descModel,
		TextArea:      textModel,
		selectedInput: m.selectedInput,
		styles:        m.styles,
		notifications: m.notifications,
		Keys:          m.Keys,
	}, tea.Batch(titleCmd, descCmd, textCmd)
}
