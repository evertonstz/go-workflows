package shared

import "github.com/evertonstz/go-workflows/models"

type (
	DidSetCurrentItemMsg struct {
		Item models.Item
	}

	SaveCommandMsg struct {
		Command string
	}

	AddNewItemMsg struct {
		Title       string
		Description string
		Command     string
	}

	CopiedToClipboardMsg struct{}

	ErrorMsg struct {
		Err error
	}
)
