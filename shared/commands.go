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

func SetCurrentItemCmd(i models.Item) tea.Cmd {
	return func() tea.Msg {
		return DidSetCurrentItemMsg{Item: i}
	}
}

func DeleteCurrentItemCmd(i int) tea.Cmd {
	return func() tea.Msg {
		return DidDeleteItemMsg{Index: i}
	}
}

func AddNewItemCmd(title, description, command string) tea.Cmd {
	return func() tea.Msg {
		return DidAddNewItemMsg{
			Title:       title,
			Description: description,
			CommandText: command,
		}
	}
}

func CloseAddNewScreenCmd() tea.Cmd {
	return func() tea.Msg {
		return DidCloseAddNewScreenMsg{}
	}
}

func CloseConfirmationModalCmd() tea.Cmd {
	return func() tea.Msg {
		return DidCloseConfirmationModalMsg{}
	}
}

func UpdateItemCmd(i models.Item) tea.Cmd {
	return func() tea.Msg {
		return DidUpdateItemMsg{Item: i}
	}
}

func SetCurrentFolderCmd(f models.FolderV2) tea.Cmd {
	return func() tea.Msg {
		return DidSetCurrentFolderMsg{Folder: f}
	}
}

func NavigatedToFolderCmd(path string) tea.Cmd {
	return func() tea.Msg {
		return DidNavigateToFolderMsg{Path: path}
	}
}

func ErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg{Err: err}
	}
}
