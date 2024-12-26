package persist

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/models"
)

type (
	Paths struct {
		DataFile   string
	}

	ErrorMsg struct {
		Err error
	}
)

type ConfigLoadedMsg struct {
	Items models.Items
}

// Comando para inicializar o gerenciador
func InitConfigManager(appName string) tea.Cmd {
	return func() tea.Msg {
		// Obtém os caminhos usando a biblioteca xdg
		dataFile, err := xdg.DataFile(fmt.Sprintf("%s/data.json", appName))
		if err != nil {
			return ErrorMsg{Err: err}
		}
		
		// Cria os diretórios se não existirem
		err = os.MkdirAll(xdg.ConfigHome+"/"+appName, os.ModePerm)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return Paths{
			DataFile:   dataFile,
		}
	}
}

// Comando para carregar e analisar o arquivo de configuração JSON.
func LoadConfigFile(path string) tea.Cmd {
	return func() tea.Msg {
		data, err := os.ReadFile(path)
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("falha ao ler o arquivo: %w", err)}
		}

		var config models.Items
		if err := json.Unmarshal(data, &config); err != nil {
			return ErrorMsg{Err: fmt.Errorf("falha ao analisar JSON: %w", err)}
		}

		return ConfigLoadedMsg{Items: config}
	}}

// comando para salvar o arquivo de configuração JSON.
func SaveConfigFile(path string, data models.Items) tea.Cmd {
	return func() tea.Msg {
		config, err := json.Marshal(data)
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("falha ao serializar JSON: %w", err)}
		}

		if err := os.WriteFile(path, config, 0644); err != nil {
			return ErrorMsg{Err: fmt.Errorf("falha ao salvar o arquivo: %w", err)}
		}

		return nil
}}