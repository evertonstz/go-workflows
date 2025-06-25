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

	// V2 specific messages
	LoadedDataFileV2Msg struct {
		Database models.DatabaseV2
	}

	PersistedFileMsg struct{}

	PersistedFileV2Msg struct{}

	// Migration messages
	MigrationCompletedMsg struct {
		FromVersion string
		ToVersion   string
	}

	DatabaseVersionMsg struct {
		Version string
	}
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

func LoadDataFileV2Cmd() tea.Cmd {
	return func() tea.Msg {
		persistenceService := di.GetService[*services.PersistenceService](di.PersistenceServiceKey)

		database, err := persistenceService.LoadDataV2()
		if err != nil {
			return shared.ErrorMsg{Err: err}
		}

		return LoadedDataFileV2Msg{Database: database}
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

func PersistListDataV2Cmd(data models.DatabaseV2) tea.Cmd {
	return func() tea.Msg {
		persistenceService := di.GetService[*services.PersistenceService](di.PersistenceServiceKey)

		if err := persistenceService.SaveDataV2(data); err != nil {
			return shared.ErrorMsg{Err: err}
		}

		return PersistedFileV2Msg{}
	}
}

func MigrateToV2Cmd() tea.Cmd {
	return func() tea.Msg {
		persistenceService := di.GetService[*services.PersistenceService](di.PersistenceServiceKey)

		// Get current version
		currentVersion, err := persistenceService.GetDatabaseVersion()
		if err != nil {
			return shared.ErrorMsg{Err: err}
		}

		// Perform migration
		if err := persistenceService.MigrateToV2(); err != nil {
			return shared.ErrorMsg{Err: err}
		}

		return MigrationCompletedMsg{
			FromVersion: currentVersion,
			ToVersion:   "2.0",
		}
	}
}

func GetDatabaseVersionCmd() tea.Cmd {
	return func() tea.Msg {
		persistenceService := di.GetService[*services.PersistenceService](di.PersistenceServiceKey)

		version, err := persistenceService.GetDatabaseVersion()
		if err != nil {
			return shared.ErrorMsg{Err: err}
		}

		return DatabaseVersionMsg{Version: version}
	}
}
