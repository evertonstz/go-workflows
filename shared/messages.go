package shared

import "github.com/evertonstz/go-workflows/models"

type (
	DidSetCurrentItemMsg struct {
		Item models.Item
	}

	DidUpdateItemMsg struct {
		Item models.Item
	}

	AddNewItemMsg struct {
		Title       string
		Description string
		Command     string
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
