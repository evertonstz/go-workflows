package main

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	confirmationmodal "github.com/evertonstz/go-workflows/components/confirmation_modal"
	"github.com/evertonstz/go-workflows/components/notification"
	"github.com/evertonstz/go-workflows/components/persist"
	addnew "github.com/evertonstz/go-workflows/screens/add_new"
	commandlist "github.com/evertonstz/go-workflows/screens/command_list"
	// "github.com/nicksnyder/go-i18n/v2/i18n" // No longer directly used here
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
	return persist.InitPersistionManagerCmd("go-workflows")
}

func new() model {
	// Ensure localizer is available. It's initialized in main.go
	if localizer == nil {
		// This is a fallback, should not happen in normal execution
		localizer = GetLocalizer("en")
	}

	persist.SetLocalizer(localizer)

	return model{
		confirmationModal: confirmationmodal.NewConfirmationModal(
			"confirmation_modal_default_message",
			"confirm_button_label",
			"cancel_button_label",
			nil, nil, localizer,
		),
		help:         help.New(),
		addNewScreen: addnew.New(localizer),
		listScreen:   commandlist.New(localizer),
		notification: notification.New(
			"app_title", // Using app_title as a default notification message for now
			false,       // false means app_title is a message ID
			localizer,
		),
		panelsStyle: panelsStyle{
			helpPanelStyle:         helpPanelStyle,
			notificationPanelStyle: notificationPanelStyle,
		},
		currentHelpHeight: 0,
		screenState:       newList,
	}
}
