package messages

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/evertonstz/go-workflows/shared/di/services"
)

// Message types for Bubble Tea integration
type (
	InitiatedPersistionMsg struct {
		DataFile string
	}

	LoadedDataFileMsg struct {
		Items models.Items
	}

	PersistedFileMsg struct{}
)

// InitPersistenceManagerCmd returns a command to initialize the persistence manager
func InitPersistenceManagerCmd() tea.Cmd {
	return func() tea.Msg {
		persistenceService := di.GetService[*services.PersistenceService](di.PersistenceServiceKey)
		return InitiatedPersistionMsg{
			DataFile: persistenceService.GetDataFilePath(),
		}
	}
}

// LoadDataFileCmd loads data from the JSON file
func LoadDataFileCmd() tea.Cmd {
	return func() tea.Msg {
		persistenceService := di.GetService[*services.PersistenceService](di.PersistenceServiceKey)

		items, err := persistenceService.LoadData()
		if err != nil {
			return shared.ErrorMsg{Err: err}
		}

		return LoadedDataFileMsg{Items: items}
	}
}

// PersistListDataCmd saves data to the JSON file
func PersistListDataCmd(data models.Items) tea.Cmd {
	return func() tea.Msg {
		persistenceService := di.GetService[*services.PersistenceService](di.PersistenceServiceKey)

		if err := persistenceService.SaveData(data); err != nil {
			return shared.ErrorMsg{Err: err}
		}

		return PersistedFileMsg{}
	}
}
