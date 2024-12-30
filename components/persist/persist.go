package persist

import (
	"encoding/json"
	"fmt"
	"os"

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
)

func InitPersistionManagerCmd(appName string) tea.Cmd {
	return func() tea.Msg {
		dataFile, err := xdg.DataFile(fmt.Sprintf("%s/data.json", appName))
		if err != nil {
			panic(err)
		}

		err = os.MkdirAll(xdg.ConfigHome+"/"+appName, os.ModePerm)
		if err != nil {
			panic(err)
		}

		return InitiatedPersistion{
			DataFile: dataFile,
		}
	}
}

func LoadDataFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		data, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		var config models.Items
		if err := json.Unmarshal(data, &config); err != nil {
			panic(err)
		}

		return LoadedDataFileMsg{Items: config}
	}
}

func SaveConfigFile(path string, data models.Items) tea.Cmd {
	return func() tea.Msg {
		config, err := json.Marshal(data)
		if err != nil {
			return shared.ErrorMsg{Err: fmt.Errorf("failed to analyze JSON: %w", err)}
		}

		if err := os.WriteFile(path, config, 0644); err != nil {
			return shared.ErrorMsg{Err: fmt.Errorf("failed saving file: %w", err)}
		}

		return nil
	}
}
