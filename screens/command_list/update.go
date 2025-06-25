package commandlist

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	confirmationmodal "github.com/evertonstz/go-workflows/components/confirmation_modal"
	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
	"github.com/evertonstz/go-workflows/shared"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case shared.DidCloseConfirmationModalMsg:
		m.currentRightPanel = textArea
	case shared.DidNavigateToFolderMsg:
		// Update current folder display
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, helpkeys.LisKeys.Esc):
			if m.currentRightPanel == modal {
				m.currentRightPanel = textArea
				return m, nil
			}
			// If at root folder, let parent handle ESC (quit app)
			// If not at root, navigable list will handle going back
		case key.Matches(msg, helpkeys.LisKeys.Delete):
			currentItem := m.navigableList.CurrentItem()
			if currentItem != nil && !currentItem.IsFolder() {
				m.showDeleteModal()
			}
		}
	}

	// Update navigable list
	m.navigableList, cmd = m.navigableList.Update(msg)
	cmds = append(cmds, cmd)

	// Update text area
	taModel, cmd := m.textArea.Update(msg)
	m.textArea = taModel.(textarea.Model)
	cmds = append(cmds, cmd)

	// Update confirmation modal if active
	if m.currentRightPanel == modal {
		confirmationModalModel, cmd := m.confirmationModal.Update(msg)
		m.confirmationModal = confirmationModalModel.(confirmationmodal.Model)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
