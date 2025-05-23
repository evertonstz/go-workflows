package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath" // Added for MkdirAll robustly

	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
)

type (
	InitiatedPersistion struct {
		DataFile string
	}

	LoadedDataFileMsg struct {
		Items models.Items
	}

	PersistedFileMsg struct{}
)

// Making XDG functions mockable
var XDGDataFile = xdg.DataFile
var OSMkdirAll = os.MkdirAll
var OSReadFile = os.ReadFile
var OSCreate = os.Create
var OSWriteFile = os.WriteFile
var JSONUnmarshal = json.Unmarshal
var JSONMarshal = json.Marshal


func InitPersistionManagerCmd(appName string) tea.Cmd {
	return func() tea.Msg {
		dataFile, err := XDGDataFile(fmt.Sprintf("%s/data.json", appName))
		if err != nil {
            return shared.ErrorMsg{Err: fmt.Errorf("xdg.DataFile failed: %w", err)}
		}

        dataDir := filepath.Dir(dataFile)
		err = OSMkdirAll(dataDir, os.ModePerm) // Use dataDir
		if err != nil {
            return shared.ErrorMsg{Err: fmt.Errorf("os.MkdirAll failed for %s: %w", dataDir, err)}
		}

		return InitiatedPersistion{
			DataFile: dataFile,
		}
	}
}

func LoadDataFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		data, err := OSReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				file, createErr := OSCreate(path)
				if createErr != nil {
                    return shared.ErrorMsg{Err: fmt.Errorf("os.Create failed for %s: %w", path, createErr)}
				}
				file.Close() 
                return LoadedDataFileMsg{Items: models.Items{Items: []models.Item{}}}
			} else {
                return shared.ErrorMsg{Err: fmt.Errorf("os.ReadFile failed for %s: %w", path, err)}
			}
		}

		if len(data) == 0 {
			return LoadedDataFileMsg{Items: models.Items{Items: []models.Item{}}}
		}

		var config models.Items
		if err := JSONUnmarshal(data, &config); err != nil {
            return shared.ErrorMsg{Err: fmt.Errorf("json.Unmarshal failed: %w", err)}
		}

		// Ensure Items slice is not nil if it's empty after unmarshalling
		if config.Items == nil {
			config.Items = []models.Item{}
		}
		return LoadedDataFileMsg{Items: config}
	}
}

// PersistDataFunc is used for allowing mocking in tests
var PersistDataFunc = PersistListData

func PersistListData(path string, data models.Items) tea.Cmd {
	return func() tea.Msg {
		config, err := JSONMarshal(data)
		if err != nil {
			return shared.ErrorMsg{Err: fmt.Errorf("failed to marshal JSON: %w", err)}
		}

		if err := OSWriteFile(path, config, 0644); err != nil {
			return shared.ErrorMsg{Err: fmt.Errorf("failed saving file %s: %w", path, err)}
		}

		return PersistedFileMsg{}
	}
}
