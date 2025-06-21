package shared

import (
	"github.com/evertonstz/go-workflows/components/persist"
	"github.com/evertonstz/go-workflows/models"
	"github.com/samber/mo"
)

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

	DidCloseAddNewScreenMsg struct{}

	CopiedToClipboardMsg struct{}

	ErrorMsg struct {
		Err error
	}

	// New Result-based messages for persistence operations
	InitiatedPersistionResultMsg struct {
		Result mo.Result[persist.InitiatedPersistion]
	}
	LoadedDataFileResultMsg struct {
		Result mo.Result[models.Items]
	}
	PersistedFileResultMsg struct {
		Result mo.Result[struct{}]
	}
)
