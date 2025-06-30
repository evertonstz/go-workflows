package textarea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
)

var (
	highlightedTextStyle = lipgloss.NewStyle()
	dateCellStyle        = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{
			Light: "#909090",
			Dark:  "#626262",
		}).
		PaddingRight(2)
)
var dateContainerStyle = lipgloss.NewStyle().Align(lipgloss.Right)

type Model struct {
	TextArea        textarea.Model
	highlightedText string
	currentItem     models.Item
	currentFolder   *models.FolderV2
	editing         bool
	err             error
}

func (m *Model) SetEditing(editing bool) {
	switch editing {
	case true:
		m.TextArea.Focus()
	case false:
		m.TextArea.Blur()
	}
	m.editing = editing
}

func New() Model {
	ti := textarea.New()
	ti.ShowLineNumbers = false
	ti.Placeholder = "Paste or type your command here..."
	ti.Prompt = ""

	return Model{
		TextArea:        ti,
		highlightedText: ti.Placeholder,
		currentFolder:   nil,
		editing:         false,
		err:             nil,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SetSize(width, height int) {
	m.TextArea.SetWidth(width)
	m.TextArea.SetHeight(height)
}

func (m *Model) SetCurrentFolder(folder models.FolderV2) {
	m.currentFolder = &folder
	m.currentItem = models.Item{} // Clear item when folder is set
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case shared.DidSetCurrentItemMsg:
		m.currentItem = msg.Item
		m.currentFolder = nil // Clear folder when item is set
		m.TextArea.SetValue(m.currentItem.Command)
	case shared.DidSetCurrentFolderMsg:
		m.currentFolder = &msg.Folder
		m.currentItem = models.Item{}
	}

	m.TextArea, cmd = m.TextArea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// TODO: remove legacy code
	var lastUpdated string
	var dateAdded string

	if m.currentFolder != nil {
		dateAdded = fmt.Sprintf("Added: %s", shared.HumanizeTime(m.currentFolder.DateAdded))
		if m.currentFolder.DateAdded.Equal(m.currentFolder.DateUpdated) {
			lastUpdated = "Updated: never"
		} else {
			lastUpdated = fmt.Sprintf("Updated: %s", shared.HumanizeTime(m.currentFolder.DateUpdated))
		}
	} else {
		dateAdded = fmt.Sprintf("Added: %s", shared.HumanizeTime(m.currentItem.DateAdded))
		if m.currentItem.DateAdded.Equal(m.currentItem.DateUpdated) {
			lastUpdated = "Updated: never"
		} else {
			lastUpdated = fmt.Sprintf("Updated: %s", shared.HumanizeTime(m.currentItem.DateUpdated))
		}
	}

	if m.editing {
		return m.TextArea.View()
	}

	rawText := m.TextArea.Value()
	highlightedText := SyntaxHighlight(rawText)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		highlightedTextStyle.
			Width(m.TextArea.Width()).
			Height(m.TextArea.Height()-2).
			Render(highlightedText),
		dateContainerStyle.
			Width(m.TextArea.Width()).
			Height(2).
			Render(lipgloss.JoinVertical(
				lipgloss.Right,
				dateCellStyle.Render(dateAdded),
				dateCellStyle.Render(lastUpdated))),
	)
}
