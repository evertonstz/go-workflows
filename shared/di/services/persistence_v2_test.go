package services

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/evertonstz/go-workflows/models"
)

func TestPersistenceService_SaveAndLoadDataV2(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_data_v2.json")

	service := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	// Create test database v2
	db := models.NewDatabaseV2()

	// Add folder
	folder := models.FolderV2{
		Name:        "Scripts",
		Description: "Development scripts",
		Path:        "/dev/scripts",
		ParentPath:  "/dev",
		DateAdded:   testTime,
		DateUpdated: testTime,
	}
	err := db.AddFolder(folder)
	if err != nil {
		t.Fatalf("Failed to add folder: %v", err)
	}

	// Add items
	item1 := models.ItemV2{
		Title:       "Test Workflow 1",
		Desc:        "Description 1",
		Command:     "echo hello",
		FolderPath:  "/dev/scripts",
		DateAdded:   testTime,
		DateUpdated: testTime,
		Tags:        []string{"test", "script"},
		Metadata:    map[string]string{"author": "test"},
	}
	item2 := models.ItemV2{
		Title:       "Test Workflow 2",
		Desc:        "Description 2",
		Command:     "ls -la",
		FolderPath:  "/",
		DateAdded:   testTime.Add(time.Hour),
		DateUpdated: testTime.Add(time.Hour),
		Tags:        []string{"utility"},
	}

	err = db.AddItem(item1)
	if err != nil {
		t.Fatalf("Failed to add item1: %v", err)
	}
	err = db.AddItem(item2)
	if err != nil {
		t.Fatalf("Failed to add item2: %v", err)
	}

	// Save v2 data
	err = service.SaveDataV2(db)
	if err != nil {
		t.Fatalf("Failed to save v2 data: %v", err)
	}

	if _, err := os.Stat(testDataFile); os.IsNotExist(err) {
		t.Fatal("Data file was not created")
	}

	// Load v2 data
	loadedDb, err := service.LoadDataV2()
	if err != nil {
		t.Fatalf("Failed to load v2 data: %v", err)
	}

	// Verify version
	if loadedDb.Version != "2.0" {
		t.Errorf("Expected version '2.0', got %q", loadedDb.Version)
	}

	// Verify folders
	if len(loadedDb.Folders) != 1 {
		t.Fatalf("Expected 1 folder, got %d", len(loadedDb.Folders))
	}

	loadedFolder := loadedDb.Folders[0]
	if loadedFolder.Name != folder.Name {
		t.Errorf("Expected folder name %q, got %q", folder.Name, loadedFolder.Name)
	}
	if loadedFolder.Path != folder.Path {
		t.Errorf("Expected folder path %q, got %q", folder.Path, loadedFolder.Path)
	}

	// Verify items
	if len(loadedDb.Items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(loadedDb.Items))
	}

	// Find items by title (order might be different)
	var loadedItem1, loadedItem2 *models.ItemV2
	for _, item := range loadedDb.Items {
		if item.Title == "Test Workflow 1" {
			loadedItem1 = &item
		} else if item.Title == "Test Workflow 2" {
			loadedItem2 = &item
		}
	}

	if loadedItem1 == nil {
		t.Fatal("Test Workflow 1 not found")
	}
	if loadedItem2 == nil {
		t.Fatal("Test Workflow 2 not found")
	}

	// Verify item1 details
	if loadedItem1.Command != item1.Command {
		t.Errorf("Expected command %q, got %q", item1.Command, loadedItem1.Command)
	}
	if loadedItem1.FolderPath != item1.FolderPath {
		t.Errorf("Expected folder path %q, got %q", item1.FolderPath, loadedItem1.FolderPath)
	}
	if len(loadedItem1.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(loadedItem1.Tags))
	}
	if loadedItem1.Metadata["author"] != "test" {
		t.Errorf("Expected metadata author 'test', got %q", loadedItem1.Metadata["author"])
	}
}

func TestPersistenceService_DetectDatabaseVersion(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_version.json")

	service := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	tests := []struct {
		name        string
		content     string
		expectedVer string
		expectError bool
	}{
		{
			name:        "v2 database",
			content:     `{"version": "2.0", "folders": [], "items": []}`,
			expectedVer: "2.0",
			expectError: false,
		},
		{
			name:        "v1 database",
			content:     `{"Items": []}`,
			expectedVer: "1.0",
			expectError: false,
		},
		{
			name:        "empty file",
			content:     "",
			expectedVer: "",
			expectError: false,
		},
		{
			name:        "no version field",
			content:     `{"items": []}`,
			expectedVer: "1.0",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.WriteFile(testDataFile, []byte(tt.content), 0o644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			data, err := os.ReadFile(testDataFile)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}

			version, err := service.detectDatabaseVersion(data)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if version != tt.expectedVer {
				t.Errorf("Expected version %q, got %q", tt.expectedVer, version)
			}
		})
	}
}

func TestPersistenceService_LoadDataBackwardCompatibility(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_compat.json")

	service := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	// Create v1 format data
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	v1Data := models.Items{
		Items: []models.Item{
			{
				Title:       "V1 Workflow",
				Desc:        "V1 Description",
				Command:     "echo v1",
				DateAdded:   testTime,
				DateUpdated: testTime,
			},
		},
	}

	// Save as v1 format
	err := service.SaveData(v1Data)
	if err != nil {
		t.Fatalf("Failed to save v1 data: %v", err)
	}

	// Load using v1 method (should work)
	loadedV1, err := service.LoadData()
	if err != nil {
		t.Fatalf("Failed to load v1 data: %v", err)
	}

	if len(loadedV1.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(loadedV1.Items))
	}

	// Load using v2 method (should migrate)
	loadedV2, err := service.LoadDataV2()
	if err != nil {
		t.Fatalf("Failed to load v1 data as v2: %v", err)
	}

	if loadedV2.Version != "2.0" {
		t.Errorf("Expected version '2.0', got %q", loadedV2.Version)
	}

	if len(loadedV2.Items) != 1 {
		t.Errorf("Expected 1 item after migration, got %d", len(loadedV2.Items))
	}

	item := loadedV2.Items[0]
	if item.Title != "V1 Workflow" {
		t.Errorf("Expected title 'V1 Workflow', got %q", item.Title)
	}
	if item.FolderPath != "/" {
		t.Errorf("Expected folder path '/', got %q", item.FolderPath)
	}
}

func TestPersistenceService_MigrateToV2(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_migrate.json")

	service := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	// Create v1 format data
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	v1Data := models.Items{
		Items: []models.Item{
			{
				Title:       "Migration Test",
				Desc:        "Test migration from v1 to v2",
				Command:     "echo migrate",
				DateAdded:   testTime,
				DateUpdated: testTime,
			},
		},
	}

	// Save as v1 format
	err := service.SaveData(v1Data)
	if err != nil {
		t.Fatalf("Failed to save v1 data: %v", err)
	}

	// Migrate to v2
	err = service.MigrateToV2()
	if err != nil {
		t.Fatalf("Failed to migrate to v2: %v", err)
	}

	// Check that backup was created
	backupPath := testDataFile + ".v1.backup"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Error("Expected backup file to be created")
	}

	// Verify that the main file is now v2
	version, err := service.GetDatabaseVersion()
	if err != nil {
		t.Fatalf("Failed to get database version: %v", err)
	}

	if version != "2.0" {
		t.Errorf("Expected version '2.0' after migration, got %q", version)
	}

	// Load and verify migrated data
	loadedV2, err := service.LoadDataV2()
	if err != nil {
		t.Fatalf("Failed to load migrated data: %v", err)
	}

	if len(loadedV2.Items) != 1 {
		t.Errorf("Expected 1 item after migration, got %d", len(loadedV2.Items))
	}

	item := loadedV2.Items[0]
	if item.Title != "Migration Test" {
		t.Errorf("Expected title 'Migration Test', got %q", item.Title)
	}
}

func TestPersistenceService_MigrateAlreadyV2(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_already_v2.json")

	service := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	// Create v2 format data
	db := models.NewDatabaseV2()
	err := service.SaveDataV2(db)
	if err != nil {
		t.Fatalf("Failed to save v2 data: %v", err)
	}

	// Try to migrate (should be no-op)
	err = service.MigrateToV2()
	if err != nil {
		t.Fatalf("Failed to migrate already v2 data: %v", err)
	}

	// Check that backup was NOT created
	backupPath := testDataFile + ".v1.backup"
	if _, err := os.Stat(backupPath); !os.IsNotExist(err) {
		t.Error("Backup file should not be created for already v2 data")
	}
}

func TestPersistenceService_GetDatabaseVersion(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_get_version.json")

	service := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	// Test non-existent file
	version, err := service.GetDatabaseVersion()
	if err != nil {
		t.Errorf("Unexpected error for non-existent file: %v", err)
	}
	if version != "" {
		t.Errorf("Expected empty version for non-existent file, got %q", version)
	}

	// Test v1 file
	v1Data := models.Items{Items: []models.Item{}}
	err = service.SaveData(v1Data)
	if err != nil {
		t.Fatalf("Failed to save v1 data: %v", err)
	}

	version, err = service.GetDatabaseVersion()
	if err != nil {
		t.Errorf("Unexpected error for v1 file: %v", err)
	}
	if version != "1.0" {
		t.Errorf("Expected version '1.0' for v1 file, got %q", version)
	}

	// Test v2 file
	db := models.NewDatabaseV2()
	err = service.SaveDataV2(db)
	if err != nil {
		t.Fatalf("Failed to save v2 data: %v", err)
	}

	version, err = service.GetDatabaseVersion()
	if err != nil {
		t.Errorf("Unexpected error for v2 file: %v", err)
	}
	if version != "2.0" {
		t.Errorf("Expected version '2.0' for v2 file, got %q", version)
	}
}

func TestPersistenceService_LoadDataV2_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "empty_v2.json")

	// Create empty file
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

	db, err := service.LoadDataV2()
	if err != nil {
		t.Fatalf("Failed to load from empty file: %v", err)
	}

	if db.Version != "2.0" {
		t.Errorf("Expected version '2.0' for new database, got %q", db.Version)
	}

	if len(db.Items) != 0 {
		t.Errorf("Expected empty items, got %d items", len(db.Items))
	}

	if len(db.Folders) != 0 {
		t.Errorf("Expected empty folders, got %d folders", len(db.Folders))
	}
}

func TestPersistenceService_LoadDataV2_NonexistentFile(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "nonexistent_v2.json")

	service := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	db, err := service.LoadDataV2()
	if err != nil {
		t.Fatalf("Failed to load from nonexistent file: %v", err)
	}

	if db.Version != "2.0" {
		t.Errorf("Expected version '2.0' for new database, got %q", db.Version)
	}

	if len(db.Items) != 0 {
		t.Errorf("Expected empty items, got %d items", len(db.Items))
	}

	if _, err := os.Stat(testDataFile); os.IsNotExist(err) {
		t.Error("Expected file to be created")
	}
}
