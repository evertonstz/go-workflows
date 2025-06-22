package commandlist

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	confirmationmodal "github.com/evertonstz/go-workflows/components/confirmation_modal"
	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/components/list"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
	"github.com/evertonstz/go-workflows/shared"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	var listModel tea.Model

	// Ensure m.localizer is not nil
	// This is a safeguard; New() should always set it.
	if m.localizer == nil {
		bund := i18n.NewBundle(language.English)
		m.localizer = i18n.NewLocalizer(bund, language.English.String())
	}
	commandListKeys := helpkeys.FullHelpKeys(m.localizer)

	switch msg := msg.(type) {
	case shared.DidCloseConfirmationModalMsg:
		m.currentRightPanel = textArea
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, commandListKeys.Esc):
			if m.currentRightPanel == modal {
				m.currentRightPanel = textArea
				return m, nil
			}
		case key.Matches(msg, commandListKeys.Delete):
			// These strings should be message IDs
			m.rebuildConfirmationModel(
				"confirm_delete_workflow_message", // New message ID
				"confirm_button_label",            // Re-use existing
				"cancel_button_label",             // Re-use existing
				tea.Batch(shared.DeleteCurrentItemCmd(m.list.CurrentItemIndex()), shared.CloseConfirmationModalCmd()),
				shared.CloseConfirmationModalCmd())
			m.currentRightPanel = modal
		}
	}

	listModel, cmd = m.list.Update(msg)
	m.list = listModel.(list.Model)
	cmds = append(cmds, cmd)

	taModel, cmd := m.textArea.Update(msg)
	m.textArea = taModel.(textarea.Model)
	cmds = append(cmds, cmd)

	if m.currentRightPanel == modal {
		confirmationModalModel, cmd := m.confirmationModal.Update(msg)
		m.confirmationModal = confirmationModalModel.(confirmationmodal.Model)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
