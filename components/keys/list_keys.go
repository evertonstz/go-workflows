package keys

import "github.com/charmbracelet/bubbles/key"

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

var LisKeys = ListKeyMap{
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete workflow"),
	),
	AddNewWorkflow: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add workflow"),
	),
	CopyWorkflow: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "copy workflow"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "close"),
	),
}
