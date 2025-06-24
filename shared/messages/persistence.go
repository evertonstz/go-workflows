package messages

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/evertonstz/go-workflows/shared/di/services"
)

type (
	InitiatedPersistionMsg struct {
		DataFile string
	}

	LoadedDataFileMsg struct {
		Items models.Items
	}

	PersistedFileMsg struct{}
)

func InitPersistenceManagerCmd() tea.Cmd {
	return func() tea.Msg {
		persistenceService := di.GetService[*services.PersistenceService](di.PersistenceServiceKey)
		return InitiatedPersistionMsg{
			DataFile: persistenceService.GetDataFilePath(),
		}
	}
}

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

func PersistListDataCmd(data models.Items) tea.Cmd {
	return func() tea.Msg {
		persistenceService := di.GetService[*services.PersistenceService](di.PersistenceServiceKey)

		if err := persistenceService.SaveData(data); err != nil {
			return shared.ErrorMsg{Err: err}
		}

		return PersistedFileMsg{}
	}
}
