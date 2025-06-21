package persist

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared" // Assuming new messages will be here
	"github.com/samber/mo"
)

type (
	InitiatedPersistion struct {
		DataFile string
	}

	// Original LoadedDataFileMsg and PersistedFileMsg might be removed or kept
	// depending on how the new Result messages are handled upstream.
	// For now, we assume they are distinct from the *ResultMsg types.
)

// Conceptual new messages (to be defined in shared/messages.go):
// type InitiatedPersistionResultMsg struct { Result mo.Result[InitiatedPersistion] }
// type LoadedDataFileResultMsg struct { Result mo.Result[models.Items] }
// type PersistedFileResultMsg struct { Result mo.Result[struct{}] }

func InitPersistionManagerCmd(appName string) tea.Cmd {
	return func() tea.Msg {
		dataFile, err := xdg.DataFile(fmt.Sprintf("%s/data.json", appName))
		if err != nil {
			return shared.InitiatedPersistionResultMsg{Result: mo.Err[InitiatedPersistion](fmt.Errorf("failed to get data file path: %w", err))}
		}

		err = os.MkdirAll(xdg.ConfigHome+"/"+appName, os.ModePerm)
		if err != nil {
			return shared.InitiatedPersistionResultMsg{Result: mo.Err[InitiatedPersistion](fmt.Errorf("failed to create config directory: %w", err))}
		}

		return shared.InitiatedPersistionResultMsg{Result: mo.Ok(InitiatedPersistion{
			DataFile: dataFile,
		})}
	}
}

func LoadDataFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				// Try to create the file if it doesn't exist
				file, createErr := os.Create(path)
				if createErr != nil {
					return shared.LoadedDataFileResultMsg{Result: mo.Err[models.Items](fmt.Errorf("failed to create data file: %w", createErr))}
				}
				file.Close() // Close immediately after creation
				// Return empty items as the file was just created
				return shared.LoadedDataFileResultMsg{Result: mo.Ok(models.Items{})}
			}
			return shared.LoadedDataFileResultMsg{Result: mo.Err[models.Items](fmt.Errorf("failed to read data file: %w", err))}
		}

		if len(data) == 0 {
			return shared.LoadedDataFileResultMsg{Result: mo.Ok(models.Items{})}
		}

		var config models.Items
		if err := json.Unmarshal(data, &config); err != nil {
			return shared.LoadedDataFileResultMsg{Result: mo.Err[models.Items](fmt.Errorf("failed to unmarshal data: %w", err))}
		}

		return shared.LoadedDataFileResultMsg{Result: mo.Ok(config)}
	}
}

func PersistListData(path string, data models.Items) tea.Cmd {
	return func() tea.Msg {
		configBytes, err := json.Marshal(data)
		if err != nil {
			return shared.PersistedFileResultMsg{Result: mo.Err[struct{}](fmt.Errorf("failed to marshal data for persistence: %w", err))}
		}

		if err := os.WriteFile(path, configBytes, 0644); err != nil {
			return shared.PersistedFileResultMsg{Result: mo.Err[struct{}](fmt.Errorf("failed to write data to file: %w", err))}
		}

		return shared.PersistedFileResultMsg{Result: mo.Ok(struct{}{})}
	}
}
