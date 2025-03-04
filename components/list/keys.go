package list

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up             key.Binding
	Down           key.Binding
	Left           key.Binding
	Right          key.Binding
	Help           key.Binding
	Quit           key.Binding
	CopyWorkflow   key.Binding
	AddNewWorkflow key.Binding
	Delete         key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// func (k KeyMap) FullHelp() [][]key.Binding {
// 	return [][]key.Binding{
// 		{k.Delete, k.Up, k.Down, k.CopyWorkflow},
// 		{k.AddNewWorkflow, k.Help, k.Quit},
// 	}
// }

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.AddNewWorkflow,k.Delete, k.CopyWorkflow},
		{k.Up, k.Down, k.Help, k.Quit},
	}
}

var Keys = KeyMap{
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
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}
