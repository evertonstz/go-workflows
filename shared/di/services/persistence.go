package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrg/xdg"

	"github.com/evertonstz/go-workflows/models"
)

type PersistenceService struct {
	dataFilePath string
	appName      string
}

func NewPersistenceService(appName string) (*PersistenceService, error) {
	dataFile, err := xdg.DataFile(fmt.Sprintf("%s/data.json", appName))
	if err != nil {
		return nil, fmt.Errorf("failed to determine data file path: %w", err)
	}

	err = os.MkdirAll(xdg.ConfigHome+"/"+appName, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	return &PersistenceService{
		dataFilePath: dataFile,
		appName:      appName,
	}, nil
}

func (p *PersistenceService) GetDataFilePath() string {
	return p.dataFilePath
}

func (p *PersistenceService) LoadData() (models.Items, error) {
	data, err := os.ReadFile(p.dataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			if _, createErr := os.Create(p.dataFilePath); createErr != nil {
				return models.Items{}, fmt.Errorf("failed to create data file: %w", createErr)
			}
			return models.Items{}, nil
		}
		return models.Items{}, fmt.Errorf("failed to read data file: %w", err)
	}

	if len(data) == 0 {
		return models.Items{}, nil
	}

	var items models.Items
	if err := json.Unmarshal(data, &items); err != nil {
		return models.Items{}, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	return items, nil
}

func (p *PersistenceService) SaveData(data models.Items) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(p.dataFilePath, jsonData, 0o644); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}
