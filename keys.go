package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up             key.Binding
	Down           key.Binding
	Left           key.Binding
	Right          key.Binding
	Help           key.Binding
	Quit           key.Binding
	Esc            key.Binding
	Enter          key.Binding
	CopyWorkflow   key.Binding
	AddNewWorkflow key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.CopyWorkflow},
		{k.AddNewWorkflow, k.Help, k.Quit},
	}
}

var keys = keyMap{
	AddNewWorkflow: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add workflow"),
	),
	CopyWorkflow: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "copy workflow"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "toggle help"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "close view"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "edit workflow"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}
