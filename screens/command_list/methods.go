package commandlist

import (
	"math"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/evertonstz/go-workflows/components/list"
	"github.com/evertonstz/go-workflows/shared"
)

func (m Model) GetAllItems() []list.ListItemInterface {
	return m.navigableList.AllItems()
}

func (m Model) GetCurrentPath() string {
	return m.navigableList.CurrentPath()
}

func (m Model) IsAtRoot() bool {
	return m.navigableList.IsAtRoot()
}

func (m *Model) setSizeForBigWidth(width, height int) {
	m.panelsStyle.leftPanelStyle = m.panelsStyle.leftPanelStyle.
		Width(int(math.Floor(float64(width) * leftPanelWidthPercentage))).
		Height(height)

	m.panelsStyle.rightPanelStyle = m.panelsStyle.rightPanelStyle.
		PaddingLeft(0).
		Width(width - m.panelsStyle.leftPanelStyle.GetWidth()).
		Height(height)

	leftWidthFrameSize, leftHeightFrameSize := m.panelsStyle.leftPanelStyle.GetFrameSize()
	rightWidthFrameSize, rightHeightFrameSize := m.panelsStyle.rightPanelStyle.GetFrameSize()

	leftPanelWidth := m.panelsStyle.leftPanelStyle.GetWidth() - leftWidthFrameSize
	rightPanelWidth := m.panelsStyle.rightPanelStyle.GetWidth() - rightWidthFrameSize

	m.navigableList.SetSize(leftPanelWidth, m.panelsStyle.leftPanelStyle.GetHeight()-leftHeightFrameSize)
	m.textArea.SetSize(rightPanelWidth, m.panelsStyle.rightPanelStyle.GetHeight()-rightHeightFrameSize)
}

func (m *Model) setSizeForSmallWidth(width, height int) {
	m.panelsStyle.leftPanelStyle = m.panelsStyle.leftPanelStyle.
		Width(width).
		Height(int(math.Floor(float64(height) * leftPanelWidthPercentage)))

	m.panelsStyle.rightPanelStyle = m.panelsStyle.rightPanelStyle.
		PaddingLeft(1).
		Width(width).
		Height(height - m.panelsStyle.leftPanelStyle.GetHeight())

	leftWidthFrameSize, leftHeightFrameSize := m.panelsStyle.leftPanelStyle.GetFrameSize()
	rightWidthFrameSize, rightHeightFrameSize := m.panelsStyle.rightPanelStyle.GetFrameSize()

	leftPanelWidth := m.panelsStyle.leftPanelStyle.GetWidth() - leftWidthFrameSize
	rightPanelWidth := m.panelsStyle.rightPanelStyle.GetWidth() - rightWidthFrameSize

	m.navigableList.SetSize(leftPanelWidth, m.panelsStyle.leftPanelStyle.GetHeight()-leftHeightFrameSize)
	m.textArea.SetSize(rightPanelWidth, m.panelsStyle.rightPanelStyle.GetHeight()-rightHeightFrameSize)
}

func (m *Model) SetSize(width, height int, smallWidth bool) {
	m.isSmallWidth = smallWidth
	if m.isSmallWidth {
		m.setSizeForSmallWidth(width, height)
	} else {
		m.setSizeForBigWidth(width, height)
	}
}

func (m *Model) showDeleteModal() {
	m.confirmationModal = m.deleteConfirmationModalBuilder(
		tea.Batch(shared.DeleteCurrentItemCmd(m.navigableList.CurrentItemIndex()), shared.CloseConfirmationModalCmd()),
		shared.CloseConfirmationModalCmd())
	m.currentRightPanel = modal
}

func (m *Model) InitializeDatabase() {
	if m.databaseManager != nil {
		m.navigableList.SetDatabase(m.databaseManager)
		m.loadInitialContent()
	}
}

func (m *Model) loadInitialContent() {
	currentItem := m.navigableList.CurrentItem()
	if currentItem == nil {
		return
	}

	if currentItem.IsFolder() {
		folder := currentItem.(list.FolderItem).GetFolder()
		m.textArea.SetCurrentFolder(folder)

		if m.databaseManager != nil {
			subfolders, items, err := m.databaseManager.GetFolderContents(folder.Path)
			if err != nil {
				m.textArea.TextArea.SetValue("Error loading folder contents: " + err.Error())
			} else {
				content := "📁 " + folder.Name + "\n" + folder.Description + "\n\n"
				content += "Contents:\n"
				for _, subfolder := range subfolders {
					content += "📁 " + subfolder.Name + " - " + subfolder.Description + "\n"
				}
				for _, item := range items {
					content += "📄 " + item.Title + " - " + item.Desc + "\n"
				}
				m.textArea.TextArea.SetValue(content)
			}
		}
	} else {
		workflowItem := currentItem.(list.WorkflowItem).GetItem()
		m.textArea.TextArea.SetValue(workflowItem.Command)
	}
}
