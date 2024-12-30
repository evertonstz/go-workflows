package shared

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/atotto/clipboard"
)

func CopyToClipboardCmd(t string) tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(t)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return CopiedToClipboardMsg{}
	}
}

func SaveCurrentItemCmd(c string) tea.Cmd {
	return func() tea.Msg {
		return SaveCommandMsg{Command: c}
	}
}
