package keys

import (
	"github.com/charmbracelet/bubbles/key"

	"github.com/evertonstz/go-workflows/shared/di/services"
)

type KeyBuilder struct {
	i18n *services.I18nService
}

func NewKeyBuilder(i18n *services.I18nService) *KeyBuilder {
	return &KeyBuilder{i18n: i18n}
}

type NavigationKeySet struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
}

type ActionKeySet struct {
	Submit key.Binding
	Close  key.Binding
	Help   key.Binding
	Quit   key.Binding
	Enter  key.Binding
	Esc    key.Binding
}

type WorkflowActionKeySet struct {
	AddNewWorkflow key.Binding
	Delete         key.Binding
	CopyWorkflow   key.Binding
}

func (b *KeyBuilder) Navigation() NavigationKeySet {
	return NavigationKeySet{
		Up:    b.key("up", "↑", "key_help_move_up"),
		Down:  b.key("down", "↓", "key_help_move_down"),
		Left:  b.key("left", "←", "key_help_move_left"),
		Right: b.key("right", "→", "key_help_move_right"),
	}
}

func (b *KeyBuilder) Actions() ActionKeySet {
	return ActionKeySet{
		Submit: b.key("enter", "enter", "key_help_submit"),
		Close:  b.key("esc", "esc", "key_help_close"),
		Help:   b.key("ctrl+h", "ctrl+h", "key_help_toggle_help"),
		Quit:   b.key("ctrl+c", "ctrl+c", "key_help_quit"),
		Enter:  b.key("enter", "enter", "key_help_open_folder"),
		Esc:    b.key("esc", "esc", "key_help_close"),
	}
}

func (b *KeyBuilder) WorkflowActions() WorkflowActionKeySet {
	return WorkflowActionKeySet{
		AddNewWorkflow: b.key("a", "a", "key_help_add_workflow"),
		Delete:         b.key("d", "d", "key_help_delete_workflow"),
		CopyWorkflow:   b.key("y", "y", "key_help_copy_workflow"),
	}
}

func (b *KeyBuilder) key(keys, short, helpKey string) key.Binding {
	return key.NewBinding(
		key.WithKeys(keys),
		key.WithHelp(short, b.i18n.Translate(helpKey)),
	)
}
