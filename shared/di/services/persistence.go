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

// DatabaseVersion represents the version structure in the JSON file
type DatabaseVersion struct {
	Version string `json:"version,omitempty"`
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

// detectDatabaseVersion checks the version field in the JSON file
func (p *PersistenceService) detectDatabaseVersion(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil // Empty file, no version
	}

	var version DatabaseVersion
	if err := json.Unmarshal(data, &version); err != nil {
		// If it fails to unmarshal with version field, assume v1
		return "1.0", nil
	}

	if version.Version == "" {
		return "1.0", nil // No version field, assume v1
	}

	return version.Version, nil
}

// LoadData loads data and returns the appropriate format based on version
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

	version, err := p.detectDatabaseVersion(data)
	if err != nil {
		return models.Items{}, fmt.Errorf("failed to detect database version: %w", err)
	}

	switch version {
	case "2.0":
		var dbV2 models.DatabaseV2
		if err := json.Unmarshal(data, &dbV2); err != nil {
			return models.Items{}, fmt.Errorf("failed to unmarshal v2 JSON data: %w", err)
		}
		// Convert v2 to v1 format for backward compatibility
		return dbV2.ToV1(), nil

	default: // v1.0 or no version
		var items models.Items
		if err := json.Unmarshal(data, &items); err != nil {
			return models.Items{}, fmt.Errorf("failed to unmarshal v1 JSON data: %w", err)
		}
		return items, nil
	}
}

// LoadDataV2 loads data and returns v2 format, migrating v1 if necessary
func (p *PersistenceService) LoadDataV2() (models.DatabaseV2, error) {
	data, err := os.ReadFile(p.dataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			if _, createErr := os.Create(p.dataFilePath); createErr != nil {
				return models.DatabaseV2{}, fmt.Errorf("failed to create data file: %w", createErr)
			}
			return models.NewDatabaseV2(), nil
		}
		return models.DatabaseV2{}, fmt.Errorf("failed to read data file: %w", err)
	}

	if len(data) == 0 {
		return models.NewDatabaseV2(), nil
	}

	version, err := p.detectDatabaseVersion(data)
	if err != nil {
		return models.DatabaseV2{}, fmt.Errorf("failed to detect database version: %w", err)
	}

	switch version {
	case "2.0":
		var dbV2 models.DatabaseV2
		if err := json.Unmarshal(data, &dbV2); err != nil {
			return models.DatabaseV2{}, fmt.Errorf("failed to unmarshal v2 JSON data: %w", err)
		}
		return dbV2, nil

	default: // v1.0 or no version
		var itemsV1 models.Items
		if err := json.Unmarshal(data, &itemsV1); err != nil {
			return models.DatabaseV2{}, fmt.Errorf("failed to unmarshal v1 JSON data: %w", err)
		}
		// Migrate v1 to v2
		return models.MigrateV1ToV2(itemsV1), nil
	}
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

// SaveDataV2 saves data in v2 format
func (p *PersistenceService) SaveDataV2(data models.DatabaseV2) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal v2 JSON: %w", err)
	}

	if err := os.WriteFile(p.dataFilePath, jsonData, 0o644); err != nil {
		return fmt.Errorf("failed to save v2 file: %w", err)
	}

	return nil
}

// MigrateToV2 migrates existing v1 data to v2 format
func (p *PersistenceService) MigrateToV2() error {
	// Load existing data
	data, err := os.ReadFile(p.dataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// No existing file, create new v2 database
			return p.SaveDataV2(models.NewDatabaseV2())
		}
		return fmt.Errorf("failed to read data file: %w", err)
	}

	if len(data) == 0 {
		// Empty file, create new v2 database
		return p.SaveDataV2(models.NewDatabaseV2())
	}

	version, err := p.detectDatabaseVersion(data)
	if err != nil {
		return fmt.Errorf("failed to detect database version: %w", err)
	}

	if version == "2.0" {
		return nil // Already v2, no migration needed
	}

	// Load v1 data and migrate
	var itemsV1 models.Items
	if err := json.Unmarshal(data, &itemsV1); err != nil {
		return fmt.Errorf("failed to unmarshal v1 JSON data: %w", err)
	}

	// Create backup of v1 data
	backupPath := p.dataFilePath + ".v1.backup"
	if err := os.WriteFile(backupPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}

	// Migrate and save v2 data
	dbV2 := models.MigrateV1ToV2(itemsV1)
	return p.SaveDataV2(dbV2)
}

// GetDatabaseVersion returns the version of the current database
func (p *PersistenceService) GetDatabaseVersion() (string, error) {
	data, err := os.ReadFile(p.dataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No file exists
		}
		return "", fmt.Errorf("failed to read data file: %w", err)
	}

	return p.detectDatabaseVersion(data)
}
