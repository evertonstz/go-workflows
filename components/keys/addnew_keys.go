package keys

import "github.com/charmbracelet/bubbles/key"

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

var AddNewKeys = AddNewKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "move right"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "move left"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "toggle help"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
	Close: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "close"),
	),
}
