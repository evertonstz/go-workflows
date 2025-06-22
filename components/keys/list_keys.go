package keys

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// This type is not used directly by the list.Model,
// but it's useful for defining the help view.
type ListKeyMap struct {
	Up             key.Binding
	Down           key.Binding
	Left           key.Binding
	Right          key.Binding
	Help           key.Binding
	Quit           key.Binding
	CopyWorkflow   key.Binding
	AddNewWorkflow key.Binding
	Delete         key.Binding
	Esc            key.Binding
}

func (k ListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k ListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.AddNewWorkflow, k.Delete, k.CopyWorkflow},
		{k.Up, k.Down, k.Help, k.Quit},
	}
}

var defaultKeys = list.DefaultKeyMap()

var ListKeys = func(localizer *i18n.Localizer) list.KeyMap { // Return the original list.KeyMap
	return list.KeyMap{ // Use the original list.KeyMap
		// Default list navigation
		CursorUp:             defaultKeys.CursorUp,
		CursorDown:           defaultKeys.CursorDown,
		NextPage:             defaultKeys.NextPage,
		PrevPage:             defaultKeys.PrevPage,
		GoToStart:            defaultKeys.GoToStart,
		GoToEnd:              defaultKeys.GoToEnd,
		Filter:               defaultKeys.Filter,
		ClearFilter:          defaultKeys.ClearFilter,
		CancelWhileFiltering: defaultKeys.CancelWhileFiltering,
		AcceptWhileFiltering: defaultKeys.AcceptWhileFiltering,
		ShowFullHelp:         defaultKeys.ShowFullHelp,
		CloseFullHelp:        defaultKeys.CloseFullHelp,

		// Customizing help text for existing actions
		// RemoveItem: key.NewBinding( // Commenting out due to persistent "unknown field" error
		// 	key.WithKeys("d"),
		// 	key.WithHelp("d", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_delete_workflow"})),
		// ),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"), // Default is "q", "ctrl+c"
			key.WithHelp("ctrl+c/q", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_quit"})),
		),
		// ShowFullHelp and CloseFullHelp already assigned from defaultKeys above
		// We can override them here if we need different help text for them specifically in this map
		// For example, if their help text needed localization:
		// ShowFullHelp: key.NewBinding(
		// 	key.WithKeys("?"),
		// 	key.WithHelp("?", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_toggle_help"})),
		// ),
		// CloseFullHelp: key.NewBinding(
		// 	key.WithKeys("?"),
		// 	key.WithHelp("?", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_close_help"})),
		// ),
		// Choose (Enter) is handled by the list component by default and is not a field in list.KeyMap

		// Assigning custom actions to keys that might not be in default KeyMap or overriding
		// Note: list.Model itself doesn't directly use AddNewWorkflow, CopyWorkflow, etc.
		// These are handled by the parent model based on key presses.
		// So, we define them here for the help view, but they won't directly control list.Model behavior.
		// The actual key matching for these custom actions will happen in the Update methods of the screen models.

		// For help display, we can add them as additional bindings if needed,
		// but they won't be part of the list.KeyMap that list.Model uses internally for its own operations.
		// The ListKeyMap struct (custom one) is better for full help display.
		// For now, focusing on what list.KeyMap directly uses.
	}
}

// This is our custom key map structure for displaying full help.
// It's separate from the list.KeyMap used by the component.
var FullHelpKeys = func(localizer *i18n.Localizer) ListKeyMap {
	return ListKeyMap{
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_delete_workflow"})),
		),
		AddNewWorkflow: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_add_workflow"})),
		),
		CopyWorkflow: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_copy_workflow"})),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_move_up"})),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_move_down"})),
		),
		Help: key.NewBinding(
			key.WithKeys("?"), // Changed from ctrl+h to ? to match common help patterns
			key.WithHelp("?", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_toggle_help"})),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("ctrl+c/q", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_quit"})),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_close"})),
		),
	}
}
