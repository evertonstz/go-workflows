package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
	"github.com/evertonstz/go-workflows/shared/di/services"
)

// FolderItem represents a folder in the list
type FolderItem struct {
	folder models.FolderV2
}

func (f FolderItem) Title() string              { return "📁 " + f.folder.Name }
func (f FolderItem) Description() string        { return f.folder.Description }
func (f FolderItem) FilterValue() string        { return f.folder.Name }
func (f FolderItem) IsFolder() bool             { return true }
func (f FolderItem) GetFolder() models.FolderV2 { return f.folder }

// WorkflowItem represents a workflow item in the list
type WorkflowItem struct {
	item models.ItemV2
}

func (w WorkflowItem) Title() string          { return "📄 " + w.item.Title }
func (w WorkflowItem) Description() string    { return w.item.Desc }
func (w WorkflowItem) FilterValue() string    { return w.item.Title }
func (w WorkflowItem) IsFolder() bool         { return false }
func (w WorkflowItem) GetItem() models.ItemV2 { return w.item }

// ListItemInterface defines the common interface for both folders and items
type ListItemInterface interface {
	list.Item
	IsFolder() bool
}

// NavigableModel extends the basic list model with folder navigation
type NavigableModel struct {
	list            list.Model
	currentPath     string
	lastSelectedIdx int
	database        *services.DatabaseManagerV2
}

func (m NavigableModel) CurrentItem() ListItemInterface {
	selectedItem := m.list.SelectedItem()
	if selectedItem == nil {
		return nil
	}

	item, ok := selectedItem.(ListItemInterface)
	if !ok {
		return nil
	}

	return item
}

func (m NavigableModel) CurrentItemIndex() int {
	return m.list.Index()
}

func (m NavigableModel) CurrentPath() string {
	return m.currentPath
}

func (m NavigableModel) IsAtRoot() bool {
	return m.currentPath == "/"
}

func (m NavigableModel) AllItems() []ListItemInterface {
	var items []ListItemInterface
	for _, i := range m.list.Items() {
		if item, ok := i.(ListItemInterface); ok {
			items = append(items, item)
		}
	}
	return items
}

func (m NavigableModel) Init() tea.Cmd {
	return nil
}

func (m *NavigableModel) SetSize(width, height int) {
	h, v := docStyle.GetFrameSize()
	m.list.SetSize(width-h, height-v)
}

func (m *NavigableModel) SetDatabase(db *services.DatabaseManagerV2) {
	m.database = db
	m.loadFolderContents(m.currentPath)
	// Initialize with the first item selected
	m.lastSelectedIdx = 0
}

func (m *NavigableModel) ReloadCurrentFolder() {
	m.loadFolderContents(m.currentPath)
	// Reset selection tracking to trigger right panel update
	m.lastSelectedIdx = -1
}

func (m *NavigableModel) loadFolderContents(folderPath string) {
	if m.database == nil {
		return
	}

	var listItems []list.Item

	// Get folder contents using DatabaseManagerV2
	subfolders, items, err := m.database.GetFolderContents(folderPath)
	if err != nil {
		// Handle error gracefully - could set an empty list or show error
		m.list.SetItems([]list.Item{})
		return
	}

	// Add folders first
	for _, folder := range subfolders {
		listItems = append(listItems, FolderItem{folder: folder})
	}

	// Add items
	for _, item := range items {
		listItems = append(listItems, WorkflowItem{item: item})
	}

	m.list.SetItems(listItems)
	m.list.Select(0) // Reset selection to top
}

func (m *NavigableModel) NavigateToFolder(folderPath string) tea.Cmd {
	m.currentPath = folderPath
	m.loadFolderContents(folderPath)

	// Force update of right panel with the newly selected item
	m.lastSelectedIdx = -1 // Force trigger of setCurrentItemCmd on next update

	// Immediately trigger right panel update for the first item in the new folder
	var cmds []tea.Cmd
	cmds = append(cmds, shared.NavigatedToFolderCmd(folderPath))
	cmds = m.setCurrentItemCmd(cmds)

	return tea.Batch(cmds...)
}

func (m *NavigableModel) NavigateUp() tea.Cmd {
	if m.IsAtRoot() {
		return nil // Can't go up from root
	}

	// Calculate parent path
	parentPath := "/"
	if m.currentPath != "/" {
		pathParts := []rune(m.currentPath)
		lastSlash := -1
		for i := len(pathParts) - 1; i >= 0; i-- {
			if pathParts[i] == '/' && i > 0 {
				lastSlash = i
				break
			}
		}
		if lastSlash > 0 {
			parentPath = string(pathParts[:lastSlash])
		}
	}

	return m.NavigateToFolder(parentPath)
}

func (m NavigableModel) setCurrentItemCmd(cmds []tea.Cmd) []tea.Cmd {
	currentItem := m.CurrentItem()
	if currentItem == nil {
		return cmds
	}

	if currentItem.IsFolder() {
		folder := currentItem.(FolderItem).GetFolder()
		cmds = append(cmds, shared.SetCurrentFolderCmd(folder))
	} else {
		item := currentItem.(WorkflowItem).GetItem()
		// Convert ItemV2 to Item for backward compatibility
		v1Item := models.Item{
			Title:       item.Title,
			Desc:        item.Desc,
			Command:     item.Command,
			DateAdded:   item.DateAdded,
			DateUpdated: item.DateUpdated,
		}
		cmds = append(cmds, shared.SetCurrentItemCmd(v1Item))
	}
	return cmds
}

func (m NavigableModel) Update(msg tea.Msg) (NavigableModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case shared.DidAddNewItemMsg:
		// Reload current folder to show new item
		m.loadFolderContents(m.currentPath)
		return m, nil

	case shared.DidDeleteItemMsg:
		// Remove item from list
		m.list.RemoveItem(msg.Index)
		if m.list.Index() >= len(m.list.Items()) {
			newIndex := len(m.list.Items()) - 1
			if newIndex < 0 {
				newIndex = 0
			}
			m.list.Select(newIndex)
		}
		return m, nil

	case shared.DidUpdateItemMsg:
		// Reload current folder to show updated item
		m.loadFolderContents(m.currentPath)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, helpkeys.LisKeys.CopyWorkflow):
			currentItem := m.CurrentItem()
			if currentItem != nil && !currentItem.IsFolder() {
				workflowItem := currentItem.(WorkflowItem)
				return m, shared.CopyToClipboardCmd(workflowItem.GetItem().Command)
			}

		case key.Matches(msg, helpkeys.LisKeys.Enter):
			currentItem := m.CurrentItem()
			if currentItem != nil && currentItem.IsFolder() {
				folderItem := currentItem.(FolderItem)
				cmd := m.NavigateToFolder(folderItem.GetFolder().Path)
				return m, cmd
			}
			// For items, do nothing (like current behavior)

		case key.Matches(msg, helpkeys.LisKeys.Esc):
			// Navigate up if not at root
			if !m.IsAtRoot() {
				cmd := m.NavigateUp()
				return m, cmd
			}
			// If at root, let the parent handle ESC (quit app)
		}
	}

	var c tea.Cmd
	m.list, c = m.list.Update(msg)
	cmds = append(cmds, c)

	if m.list.Index() != m.lastSelectedIdx {
		cmds = m.setCurrentItemCmd(cmds)
		m.lastSelectedIdx = m.list.Index()
	}

	return m, tea.Batch(cmds...)
}

func (m NavigableModel) View() string {
	return docStyle.Render(m.list.View())
}

func NewNavigable() NavigableModel {
	delegate := list.NewDefaultDelegate()

	// Customize the delegate to show folders and items differently
	delegate.ShowDescription = true

	m := NavigableModel{
		list:        list.New([]list.Item{}, delegate, 0, 0),
		currentPath: "/",
	}

	m.list.SetShowTitle(false)
	m.list.SetShowHelp(false)
	m.Init()

	return m
}
