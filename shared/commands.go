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
			// TODO: Emit an error message to the future notification system
			return ErrorMsg{Err: err}
		}
		return CopiedToClipboardMsg{}
	}
}

func SetCurrentItemCmd(i models.Item) tea.Cmd {
	return func() tea.Msg {
		return DidSetCurrentItemMsg{Item: i}
	}
}
