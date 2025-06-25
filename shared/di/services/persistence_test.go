package services

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/evertonstz/go-workflows/models"
)

func TestPersistenceService_SaveAndLoadData(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_data.json")

	service := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	testItems := models.Items{
		Items: []models.Item{
			{
				Title:       "Test Workflow 1",
				Desc:        "Description 1",
				Command:     "echo hello",
				DateAdded:   testTime,
				DateUpdated: testTime,
			},
			{
				Title:       "Test Workflow 2",
				Desc:        "Description 2",
				Command:     "ls -la",
				DateAdded:   testTime.Add(time.Hour),
				DateUpdated: testTime.Add(time.Hour),
			},
		},
	}

	err := service.SaveData(testItems)
	if err != nil {
		t.Fatalf("Failed to save data: %v", err)
	}

	if _, err := os.Stat(testDataFile); os.IsNotExist(err) {
		t.Fatal("Data file was not created")
	}

	loadedItems, err := service.LoadData()
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	if len(loadedItems.Items) != len(testItems.Items) {
		t.Fatalf("Expected %d items, got %d", len(testItems.Items), len(loadedItems.Items))
	}

	for i, expected := range testItems.Items {
		actual := loadedItems.Items[i]
		if actual.Title != expected.Title {
			t.Errorf("Item %d: expected title %q, got %q", i, expected.Title, actual.Title)
		}
		if actual.Desc != expected.Desc {
			t.Errorf("Item %d: expected desc %q, got %q", i, expected.Desc, actual.Desc)
		}
		if actual.Command != expected.Command {
			t.Errorf("Item %d: expected command %q, got %q", i, expected.Command, actual.Command)
		}
		if !actual.DateAdded.Equal(expected.DateAdded) {
			t.Errorf("Item %d: expected date added %v, got %v", i, expected.DateAdded, actual.DateAdded)
		}
		if !actual.DateUpdated.Equal(expected.DateUpdated) {
			t.Errorf("Item %d: expected date updated %v, got %v", i, expected.DateUpdated, actual.DateUpdated)
		}
	}
}

func TestPersistenceService_LoadData_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "empty_data.json")

	file, err := os.Create(testDataFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("Failed to close test file: %v", err)
	}

	service := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	items, err := service.LoadData()
	if err != nil {
		t.Fatalf("Failed to load from empty file: %v", err)
	}

	if len(items.Items) != 0 {
		t.Errorf("Expected empty items, got %d items", len(items.Items))
	}
}

func TestPersistenceService_LoadData_NonexistentFile(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "nonexistent.json")

	service := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	items, err := service.LoadData()
	if err != nil {
		t.Fatalf("Failed to load from nonexistent file: %v", err)
	}

	if len(items.Items) != 0 {
		t.Errorf("Expected empty items, got %d items", len(items.Items))
	}

	if _, err := os.Stat(testDataFile); os.IsNotExist(err) {
		t.Error("Expected file to be created")
	}
}

func TestPersistenceService_SaveData_InvalidData(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_data.json")

	service := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	validItems := models.Items{
		Items: []models.Item{
			{
				Title:   "Valid Item",
				Desc:    "Valid Description",
				Command: "echo test",
			},
		},
	}

	err := service.SaveData(validItems)
	if err != nil {
		t.Fatalf("Failed to save valid data: %v", err)
	}
}

func TestPersistenceService_GetDataFilePath(t *testing.T) {
	expectedPath := "/test/path/data.json"
	service := &PersistenceService{
		dataFilePath: expectedPath,
		appName:      "test-app",
	}

	actualPath := service.GetDataFilePath()
	if actualPath != expectedPath {
		t.Errorf("Expected path %q, got %q", expectedPath, actualPath)
	}
}

func TestNewPersistenceService(t *testing.T) {
	appName := "test-app"
	service, err := NewPersistenceService(appName)
	if err != nil {
		t.Fatalf("Failed to create persistence service: %v", err)
	}

	if service.appName != appName {
		t.Errorf("Expected app name %q, got %q", appName, service.appName)
	}

	if service.dataFilePath == "" {
		t.Error("Expected data file path to be set")
	}
}
