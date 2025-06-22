package keys

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type AddNewKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Help   key.Binding
	Close  key.Binding
	Submit key.Binding
}

func (k AddNewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Close}
}

func (k AddNewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Help, k.Close},
	}
}

var AddNewKeys = func(localizer *i18n.Localizer) AddNewKeyMap {
	return AddNewKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_move_up"})),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_move_down"})),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_move_right"})),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_move_left"})),
		),
		Help: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_toggle_help"})),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_submit"})),
		),
		Close: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "key_help_close"})),
		),
	}
}
