package commandlist

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/components/list"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
	confirmationmodal "github.com/evertonstz/go-workflows/components/confirmation_modal"
	helpkeys "github.com/evertonstz/go-workflows/components/keys"
    "github.com/evertonstz/go-workflows/shared"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case shared.DidCloseConfirmationModalMsg:
		m.currentRightPanel = textArea
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, helpkeys.LisKeys.Delete):
			m.rebuildConfirmationModel("Are you sure you want to delete this workflow?",
				"Yes",
				"No",
				tea.Batch(shared.DeleteCurrentItemCmd(m.list.CurrentItemIndex()), shared.CloseConfirmationModalCmd()),
				shared.CloseConfirmationModalCmd())
			m.currentRightPanel = modal
	}}

	var cmd tea.Cmd
	var listModel tea.Model
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
