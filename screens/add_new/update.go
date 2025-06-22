package addnew

import (
	"github.com/charmbracelet/bubbles/key"

	tea "github.com/charmbracelet/bubbletea"
	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/components/notification"
	"github.com/evertonstz/go-workflows/shared"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.localizer == nil {
		// This should not happen if New() is always used to create the model.
		// Consider logging an error or returning a specific error message.
		// For now, let's try to proceed with a default localizer to avoid a panic.
		// This is a temporary safeguard.
		bund := i18n.NewBundle(language.English)
		m.localizer = i18n.NewLocalizer(bund, language.English.String())
	}
	addNewKeyMap := helpkeys.AddNewKeys(m.localizer)

	titleModel, titleCmd := m.Title.Update(msg)
	descModel, descCmd := m.Description.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, addNewKeyMap.Down):
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
		case key.Matches(msg, addNewKeyMap.Up):
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
		case key.Matches(msg, addNewKeyMap.Right):
			if m.selectedInput == submit {
				return m.focusInput(close)
			}
		case key.Matches(msg, addNewKeyMap.Left):
			if m.selectedInput == close {
				return m.focusInput(submit)
			}
		case key.Matches(msg, addNewKeyMap.Close):
			m.ResetForm()
			return m, shared.CloseAddNewScreenCmd()
		case key.Matches(msg, addNewKeyMap.Submit):
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
				// TODO: Add "error_fill_all_fields" to en.json
				return m, notification.ShowNotificationCmd("error_fill_all_fields", false, nil)
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
		localizer:     m.localizer, // Persist localizer
	}, tea.Batch(titleCmd, descCmd, textCmd)
}
