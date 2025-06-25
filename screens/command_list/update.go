package commandlist

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	confirmationmodal "github.com/evertonstz/go-workflows/components/confirmation_modal"
	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/components/list"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
	"github.com/evertonstz/go-workflows/shared"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case shared.DidCloseConfirmationModalMsg:
		m.currentRightPanel = textArea
	case shared.DidDeleteItemMsg:
		if m.databaseManager != nil {
			currentItem := m.navigableList.CurrentItem()
			if currentItem != nil && !currentItem.IsFolder() {
				workflowItem := currentItem.(list.WorkflowItem).GetItem()
				err := m.databaseManager.DeleteItem(workflowItem.ID)
				if err != nil {
					return m, shared.ErrorCmd(err)
				}
				m.navigableList.ReloadCurrentFolder()
			}
		}
		return m, nil
	case shared.DidAddNewItemMsg:
		if m.databaseManager != nil {
			currentPath := m.navigableList.CurrentPath()
			_, err := m.databaseManager.CreateItem(
				msg.Title,
				msg.Description,
				msg.CommandText,
				currentPath,
				[]string{},          // empty tags
				map[string]string{}, // empty metadata
			)
			if err != nil {
				return m, shared.ErrorCmd(err)
			}

			m.navigableList.ReloadCurrentFolder()
		}
		return m, nil
	case shared.DidNavigateToFolderMsg:
		return m, nil
	case shared.DidSetCurrentItemMsg:
		m.currentRightPanel = textArea
	case shared.DidSetCurrentFolderMsg:
		m.textArea.SetCurrentFolder(msg.Folder)

		if m.databaseManager != nil {
			subfolders, items, err := m.databaseManager.GetFolderContents(msg.Folder.Path)
			if err != nil {
				m.textArea.TextArea.SetValue("Error loading folder contents: " + err.Error())
			} else {
				content := "📁 " + msg.Folder.Name + "\n" + msg.Folder.Description + "\n\n"
				content += "Contents:\n"
				for _, folder := range subfolders {
					content += "📁 " + folder.Name + " - " + folder.Description + "\n"
				}
				for _, item := range items {
					content += "📄 " + item.Title + " - " + item.Desc + "\n"
				}

				m.textArea.TextArea.SetValue(content)
			}
		}
		m.currentRightPanel = textArea
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, helpkeys.LisKeys.Esc):
			if m.currentRightPanel == modal {
				m.currentRightPanel = textArea
				return m, nil
			}
		case key.Matches(msg, helpkeys.LisKeys.Delete):
			currentItem := m.navigableList.CurrentItem()
			if currentItem != nil && !currentItem.IsFolder() {
				m.showDeleteModal()
			}
		}
	}

	m.navigableList, cmd = m.navigableList.Update(msg)
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
