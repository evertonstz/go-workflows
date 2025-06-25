package shared

import "github.com/evertonstz/go-workflows/models"

type (
	DidSetCurrentItemMsg struct {
		Item models.Item
	}

	DidSetCurrentFolderMsg struct {
		Folder models.FolderV2
	}

	DidNavigateToFolderMsg struct {
		Path string
	}

	DidUpdateItemMsg struct {
		Item models.Item
	}

	DidAddNewItemMsg struct {
		Title       string
		Description string
		CommandText string
	}

	DidDeleteItemMsg struct {
		Index int
	}

	DidCloseConfirmationModalMsg struct{}

	DidCloseAddNewScreenMsg struct{}

	CopiedToClipboardMsg struct{}

	ErrorMsg struct {
		Err error
	}
)
