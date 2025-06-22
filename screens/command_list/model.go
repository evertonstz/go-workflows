package commandlist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	confirmationmodal "github.com/evertonstz/go-workflows/components/confirmation_modal"
	"github.com/evertonstz/go-workflows/components/list"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	localizer *i18n.Localizer

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
		localizer         *i18n.Localizer // Add localizer to the model struct
	}
	currentRightPanel uint
)

const (
	textArea currentRightPanel = iota
	modal
)

func (m Model) Init() tea.Cmd {
	return nil
}

func New(loc *i18n.Localizer) Model {
	localizer = loc // Keep for package-level access if still needed by some functions directly
	// It's generally better to pass localizer explicitly or via struct fields.

	listModel := list.New(loc)
	textAreaModel := textarea.New()
	confirmModal := confirmationmodal.NewConfirmationModal(
		"confirmation_modal_default_message",
		"confirm_button_label",
		"cancel_button_label",
		nil,
		nil,
		loc, // Pass the received localizer
	)

	return Model{
		list:              listModel,
		confirmationModal: confirmModal,
		textArea:          textAreaModel,
		panelsStyle: panelsStyle{
			leftPanelStyle:  leftPanelStyle,
			rightPanelStyle: rightPanelStyle,
		},
		currentRightPanel: textArea,
		isSmallWidth:      false,
		localizer:         loc, // Store the localizer in the model instance
	}
}
