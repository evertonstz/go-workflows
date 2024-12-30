package shared

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/models"

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

func SaveSelectedItemCmd(c string) tea.Cmd {
	return func() tea.Msg {
		return SaveCommandMsg{Command: c}
	}
}

func SetCurrentItemCmd(i models.Item) tea.Cmd {
	return func() tea.Msg {
		return DidSetCurrentItemMsg{Item: i}
	}
}
