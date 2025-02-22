package shared

import "github.com/evertonstz/go-workflows/models"

type (
	DidSetCurrentItemMsg struct {
		Item models.Item
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

	CopiedToClipboardMsg struct{}

	ErrorMsg struct {
		Err error
	}
)
