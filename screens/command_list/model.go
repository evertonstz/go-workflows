package commandlist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	confirmationmodal "github.com/evertonstz/go-workflows/components/confirmation_modal"
	"github.com/evertonstz/go-workflows/components/list"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
)

var (
	leftPanelWidthPercentage = 0.5
	leftPanelStyle           = lipgloss.NewStyle().
					AlignHorizontal(lipgloss.Left)
	rightPanelStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Left).
			PaddingTop(1).
			Width(15).
			Height(5)
)

type (
	panelsStyle struct {
		leftPanelStyle  lipgloss.Style
		rightPanelStyle lipgloss.Style
	}

	Model struct {
		list              list.Model
		confirmationModal confirmationmodal.Model
		textArea          textarea.Model
		panelsStyle       panelsStyle
		currentRightPanel currentRightPanel
		isSmallWidth      bool
	}
	currentRightPanel uint
)

// GetAllItems retrieves all items from the underlying list.Model.
func (m Model) GetAllItems() []list.MyItem {
	return m.list.AllItems()
}

const (
	textArea currentRightPanel = iota
	modal
)

func (m Model) Init() tea.Cmd {
	return nil
}

func New() Model {
	listModel := list.New()
	textAreaModel := textarea.New()
	confirmationmodal := confirmationmodal.NewConfirmationModal("", "", "", nil, nil)

	return Model{
		list:              listModel,
		confirmationModal: confirmationmodal,
		textArea:          textAreaModel,
		panelsStyle: panelsStyle{
			leftPanelStyle:  leftPanelStyle,
			rightPanelStyle: rightPanelStyle,
		},
		currentRightPanel: textArea,
		isSmallWidth:      false,
	}
}
