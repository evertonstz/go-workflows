package persist

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var localizer *i18n.Localizer

func SetLocalizer(loc *i18n.Localizer) {
	localizer = loc
}

type (
	InitiatedPersistion struct {
		DataFile string
	}

	LoadedDataFileMsg struct {
		Items models.Items
	}

	PersistedFileMsg struct{}
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
			if os.IsNotExist(err) {
				if _, createErr := os.Create(path); createErr != nil {
					panic(createErr)
				}
			} else {
				panic(err)
			}
		}

		if len(data) == 0 {
			return LoadedDataFileMsg{Items: models.Items{}}
		}

		var config models.Items
		if err := json.Unmarshal(data, &config); err != nil {
			panic(err)
		}

		return LoadedDataFileMsg{Items: config}
	}
}

func PersistListData(path string, data models.Items) tea.Cmd {
	return func() tea.Msg {
		config, err := json.Marshal(data)
		if err != nil {
			localizedErr := localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "error_failed_to_analyze_json",
				TemplateData: map[string]interface{}{
					"Error": err.Error(),
				},
			})
			return shared.ErrorMsg{Err: fmt.Errorf(localizedErr)}
		}

		if err := os.WriteFile(path, config, 0644); err != nil {
			localizedErr := localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "error_failed_saving_file",
				TemplateData: map[string]interface{}{
					"Error": err.Error(),
				},
			})
			return shared.ErrorMsg{Err: fmt.Errorf(localizedErr)}
		}

		return PersistedFileMsg{}
	}
}
