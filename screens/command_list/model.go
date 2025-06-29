package commandlist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	confirmationmodal "github.com/evertonstz/go-workflows/components/confirmation_modal"
	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/components/list"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/evertonstz/go-workflows/shared/di/services"
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

	confirmationModalBuilder func(confirmCmd, cancelCmd tea.Cmd) confirmationmodal.Model

	Model struct {
		navigableList                  list.NavigableModel
		confirmationModal              confirmationmodal.Model
		deleteConfirmationModalBuilder confirmationModalBuilder
		textArea                       textarea.Model
		panelsStyle                    panelsStyle
		currentRightPanel              currentRightPanel
		isSmallWidth                   bool
		databaseManager                *services.DatabaseManagerV2
		Keys                           helpkeys.ListKeyMap
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

func New() Model {
	i18n := di.GetService[*services.I18nService](di.I18nServiceKey)

	navigableListModel := list.NewNavigable()
	textAreaModel := textarea.New()
	initialModal := confirmationmodal.NewConfirmationModal("", "", "", nil, nil)
	deleteConfirmationModalBuilder := func(confirmCmd, cancelCmd tea.Cmd) confirmationmodal.Model {
		modal := confirmationmodal.NewConfirmationModal(
			i18n.Translate("confirm_delete_workflow_message"),
			i18n.Translate("yes"),
			i18n.Translate("no"),
			confirmCmd,
			cancelCmd,
		)
		return modal
	}

	persistence := di.GetService[*services.PersistenceService](di.PersistenceServiceKey)
	databaseManager, err := services.NewDatabaseManagerV2(persistence)
	if err != nil {
		databaseManager = nil
	}

	return Model{
		navigableList:                  navigableListModel,
		confirmationModal:              initialModal,
		deleteConfirmationModalBuilder: deleteConfirmationModalBuilder,
		textArea:                       textAreaModel,
		Keys:                           helpkeys.NewListKeys(i18n),
		panelsStyle: panelsStyle{
			leftPanelStyle:  leftPanelStyle,
			rightPanelStyle: rightPanelStyle,
		},
		currentRightPanel: textArea,
		isSmallWidth:      false,
		databaseManager:   databaseManager,
	}
}
