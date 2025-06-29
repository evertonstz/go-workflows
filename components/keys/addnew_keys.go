package keys

import (
	"github.com/charmbracelet/bubbles/key"

	"github.com/evertonstz/go-workflows/shared/di/services"
)

type AddNewKeyMap struct {
	NavigationKeySet
	ActionKeySet
}

func (k AddNewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Close}
}

func (k AddNewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Submit, k.Close, k.Help, k.Quit},
	}
}

func NewAddNewKeys(i18n *services.I18nService) AddNewKeyMap {
	builder := NewKeyBuilder(i18n)
	navigation := builder.Navigation()
	actions := builder.Actions()

	return AddNewKeyMap{
		NavigationKeySet: navigation,
		ActionKeySet: ActionKeySet{
			Submit: actions.Submit,
			Close:  actions.Close,
			Help:   actions.Help,
			Quit:   actions.Quit,
			// Not including Enter, Esc as they're not used in add new screen
		},
	}
}
