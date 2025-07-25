package main

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	confirmationmodal "github.com/evertonstz/go-workflows/components/confirmation_modal"
	"github.com/evertonstz/go-workflows/components/notification"
	addnew "github.com/evertonstz/go-workflows/screens/add_new"
	commandlist "github.com/evertonstz/go-workflows/screens/command_list"
	"github.com/evertonstz/go-workflows/shared/messages"
)

var (
	helpPanelStyle         = lipgloss.NewStyle().PaddingLeft(2)
	notificationPanelStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Height(1).
				AlignHorizontal(lipgloss.Left)
)

type (
	termDimensions struct {
		width  int
		height int
	}

	panelsStyle struct {
		helpPanelStyle         lipgloss.Style
		notificationPanelStyle lipgloss.Style
	}

	model struct {
		confirmationModal confirmationmodal.Model
		help              help.Model
		screenState       screenState
		addNewScreen      addnew.Model
		listScreen        commandlist.Model
		persistPath       string
		currentPath       string
		notification      notification.Model
		termDimensions    termDimensions
		currentHelpHeight int
		panelsStyle       panelsStyle
	}
	screenState uint
)

const (
	addNew screenState = iota
	newList
)

func (m model) Init() tea.Cmd {
	return messages.InitPersistenceManagerCmd()
}

func new() model {
	listScreen := commandlist.New()
	listScreen.InitializeDatabase()

	return model{
		confirmationModal: confirmationmodal.NewConfirmationModal("", "", "", nil, nil),
		help:              help.New(),
		addNewScreen:      addnew.New(),
		listScreen:        listScreen,
		currentPath:       "/",
		notification:      notification.New("Workflows"),
		panelsStyle: panelsStyle{
			helpPanelStyle:         helpPanelStyle,
			notificationPanelStyle: notificationPanelStyle,
		},
		currentHelpHeight: 0,
		screenState:       newList,
	}
}
