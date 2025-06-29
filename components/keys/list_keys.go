package keys

import (
	"github.com/charmbracelet/bubbles/key"

	"github.com/evertonstz/go-workflows/shared/di/services"
)

type ListKeyMap struct {
	NavigationKeySet
	ActionKeySet
	WorkflowActionKeySet
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

func NewListKeys(i18n *services.I18nService) ListKeyMap {
	builder := NewKeyBuilder(i18n)
	navigation := builder.Navigation()
	actions := builder.Actions()
	workflowActions := builder.WorkflowActions()

	return ListKeyMap{
		NavigationKeySet:     navigation,
		ActionKeySet:         actions,
		WorkflowActionKeySet: workflowActions,
	}
}

var LisKeys ListKeyMap

func InitializeGlobalKeys(i18n *services.I18nService) {
	LisKeys = NewListKeys(i18n)
}
