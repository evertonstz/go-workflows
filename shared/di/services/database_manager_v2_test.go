package services

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/evertonstz/go-workflows/models"
)

func createTestDatabaseManager(dataFilePath string) (*DatabaseManagerV2, error) {
	persistence := &PersistenceService{
		dataFilePath: dataFilePath,
		appName:      "test-app",
	}

	validation := NewValidationService()

	err := persistence.SaveDataV2(models.NewDatabaseV2())
	if err != nil {
		return nil, err
	}

	return NewDatabaseManagerV2(persistence, validation)
}

func TestDatabaseManagerV2_CreateFolder(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_db_manager.json")

	manager, err := createTestDatabaseManager(testDataFile)
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}

	folder, err := manager.CreateFolder("scripts", "Development scripts", "/")
	if err != nil {
		t.Fatalf("Failed to create folder: %v", err)
	}

	if folder.Name != "scripts" {
		t.Errorf("Expected folder name 'scripts', got %q", folder.Name)
	}

	if folder.Path != "/scripts" {
		t.Errorf("Expected folder path '/scripts', got %q", folder.Path)
	}

	if folder.ID == "" {
		t.Error("Expected folder ID to be generated")
	}

	subfolder, err := manager.CreateFolder("utils", "Utility scripts", "/scripts")
	if err != nil {
		t.Fatalf("Failed to create subfolder: %v", err)
	}

	if subfolder.Path != "/scripts/utils" {
		t.Errorf("Expected subfolder path '/scripts/utils', got %q", subfolder.Path)
	}

	if subfolder.ParentPath != "/scripts" {
		t.Errorf("Expected parent path '/scripts', got %q", subfolder.ParentPath)
	}

	_, err = manager.CreateFolder("invalid", "Invalid folder", "/nonexistent")
	if err == nil {
		t.Error("Expected error when creating folder with non-existent parent")
	}
}

func TestDatabaseManagerV2_CreateItem(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_db_manager_items.json")

	persistence := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	err := persistence.SaveDataV2(models.NewDatabaseV2())
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	manager, err := createTestDatabaseManager(testDataFile)
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}

	_, err = manager.CreateFolder("scripts", "Scripts folder", "/")
	if err != nil {
		t.Fatalf("Failed to create folder: %v", err)
	}

	rootItem, err := manager.CreateItem(
		"Root Script",
		"A script in root",
		"echo root",
		"/",
		[]string{"test", "root"},
		map[string]string{"author": "test"},
	)
	if err != nil {
		t.Fatalf("Failed to create root item: %v", err)
	}

	if rootItem.FolderPath != "/" {
		t.Errorf("Expected folder path '/', got %q", rootItem.FolderPath)
	}

	if len(rootItem.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(rootItem.Tags))
	}

	if rootItem.Metadata["author"] != "test" {
		t.Errorf("Expected metadata author 'test', got %q", rootItem.Metadata["author"])
	}

	subItem, err := manager.CreateItem(
		"Sub Script",
		"A script in subfolder",
		"echo sub",
		"/scripts",
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to create sub item: %v", err)
	}

	if subItem.FolderPath != "/scripts" {
		t.Errorf("Expected folder path '/scripts', got %q", subItem.FolderPath)
	}

	_, err = manager.CreateItem("Invalid", "Invalid item", "echo invalid", "/nonexistent", nil, nil)
	if err == nil {
		t.Error("Expected error when creating item in non-existent folder")
	}
}

func TestDatabaseManagerV2_GetFolderContents(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_folder_contents.json")

	persistence := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	err := persistence.SaveDataV2(models.NewDatabaseV2())
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	manager, err := createTestDatabaseManager(testDataFile)
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}

	_, err = manager.CreateFolder("dev", "Development", "/")
	if err != nil {
		t.Fatalf("Failed to create dev folder: %v", err)
	}

	_, err = manager.CreateFolder("scripts", "Scripts", "/dev")
	if err != nil {
		t.Fatalf("Failed to create scripts folder: %v", err)
	}

	_, err = manager.CreateFolder("tools", "Tools", "/dev")
	if err != nil {
		t.Fatalf("Failed to create tools folder: %v", err)
	}

	_, err = manager.CreateItem("Root Item", "Root item", "echo root", "/", nil, nil)
	if err != nil {
		t.Fatalf("Failed to create root item: %v", err)
	}

	_, err = manager.CreateItem("Dev Item", "Dev item", "echo dev", "/dev", nil, nil)
	if err != nil {
		t.Fatalf("Failed to create dev item: %v", err)
	}

	_, err = manager.CreateItem("Script Item", "Script item", "echo script", "/dev/scripts", nil, nil)
	if err != nil {
		t.Fatalf("Failed to create script item: %v", err)
	}

	subfolders, items, err := manager.GetFolderContents("/")
	if err != nil {
		t.Fatalf("Failed to get root contents: %v", err)
	}

	if len(subfolders) != 1 {
		t.Errorf("Expected 1 subfolder in root, got %d", len(subfolders))
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item in root, got %d", len(items))
	}

	subfolders, items, err = manager.GetFolderContents("/dev")
	if err != nil {
		t.Fatalf("Failed to get dev contents: %v", err)
	}

	if len(subfolders) != 2 {
		t.Errorf("Expected 2 subfolders in /dev, got %d", len(subfolders))
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item in /dev, got %d", len(items))
	}

	subfolders, items, err = manager.GetFolderContents("/dev/scripts")
	if err != nil {
		t.Fatalf("Failed to get scripts contents: %v", err)
	}

	if len(subfolders) != 0 {
		t.Errorf("Expected 0 subfolders in /dev/scripts, got %d", len(subfolders))
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item in /dev/scripts, got %d", len(items))
	}
}

func TestDatabaseManagerV2_UpdateItem(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_update_item.json")

	persistence := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	err := persistence.SaveDataV2(models.NewDatabaseV2())
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	manager, err := createTestDatabaseManager(testDataFile)
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}

	_, err = manager.CreateFolder("scripts", "Scripts", "/")
	if err != nil {
		t.Fatalf("Failed to create folder: %v", err)
	}

	item, err := manager.CreateItem("Original", "Original desc", "echo original", "/", []string{"old"}, map[string]string{"version": "1"})
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	err = manager.UpdateItem(
		item.ID,
		"Updated Title",
		"Updated description",
		"echo updated",
		"/scripts",
		[]string{"new", "updated"},
		map[string]string{"version": "2", "author": "test"},
	)
	if err != nil {
		t.Fatalf("Failed to update item: %v", err)
	}

	updatedItem, err := manager.GetItem(item.ID)
	if err != nil {
		t.Fatalf("Failed to get updated item: %v", err)
	}

	if updatedItem.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %q", updatedItem.Title)
	}

	if updatedItem.FolderPath != "/scripts" {
		t.Errorf("Expected folder path '/scripts', got %q", updatedItem.FolderPath)
	}

	if len(updatedItem.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(updatedItem.Tags))
	}

	if updatedItem.Metadata["version"] != "2" {
		t.Errorf("Expected metadata version '2', got %q", updatedItem.Metadata["version"])
	}
}

func TestDatabaseManagerV2_Search(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_search.json")

	persistence := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	err := persistence.SaveDataV2(models.NewDatabaseV2())
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	manager, err := createTestDatabaseManager(testDataFile)
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}

	_, err = manager.CreateFolder("scripts", "Scripts folder", "/")
	if err != nil {
		t.Fatalf("Failed to create folder: %v", err)
	}

	_, err = manager.CreateItem("Test Script", "A test script", "echo test", "/scripts", []string{"test", "automation"}, nil)
	if err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	_, err = manager.CreateItem("Deploy Script", "Deployment script", "deploy.sh", "/scripts", []string{"deploy", "production"}, nil)
	if err != nil {
		t.Fatalf("Failed to create deploy item: %v", err)
	}

	_, err = manager.CreateItem("Backup Tool", "Backup utility", "backup.sh", "/", []string{"backup", "maintenance"}, nil)
	if err != nil {
		t.Fatalf("Failed to create backup item: %v", err)
	}
	result := manager.Search(models.SearchCriteria{Query: "test"})
	if len(result.Items) != 1 {
		t.Errorf("Expected 1 item with 'test' query, got %d", len(result.Items))
	}

	result = manager.SearchInFolder("/scripts", "")
	if len(result.Items) != 2 {
		t.Errorf("Expected 2 items in /scripts folder, got %d", len(result.Items))
	}

	result = manager.SearchByTags([]string{"deploy"})
	if len(result.Items) != 1 {
		t.Errorf("Expected 1 item with 'deploy' tag, got %d", len(result.Items))
	}

	result = manager.Search(models.SearchCriteria{FolderPath: "/"})
	if len(result.Items) != 1 {
		t.Errorf("Expected 1 item in root folder, got %d", len(result.Items))
	}
}

func TestDatabaseManagerV2_DeleteFolder(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_delete_folder.json")

	persistence := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	err := persistence.SaveDataV2(models.NewDatabaseV2())
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	manager, err := createTestDatabaseManager(testDataFile)
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}

	_, err = manager.CreateFolder("temp", "Temporary folder", "/")
	if err != nil {
		t.Fatalf("Failed to create temp folder: %v", err)
	}

	_, err = manager.CreateFolder("subfolder", "Sub folder", "/temp")
	if err != nil {
		t.Fatalf("Failed to create subfolder: %v", err)
	}

	_, err = manager.CreateItem("Temp Item", "Temporary item", "echo temp", "/temp", nil, nil)
	if err != nil {
		t.Fatalf("Failed to create temp item: %v", err)
	}

	err = manager.DeleteFolder("/temp", false)
	if err == nil {
		t.Error("Expected error when deleting non-empty folder without force")
	}

	err = manager.DeleteFolder("/temp", true)
	if err != nil {
		t.Fatalf("Failed to delete folder with force: %v", err)
	}

	_, err = manager.GetFolder("/temp")
	if err == nil {
		t.Error("Expected error when getting deleted folder")
	}

	db := manager.GetDatabase()
	for _, item := range db.Items {
		if item.FolderPath == "/temp" {
			t.Error("Found item in deleted folder")
		}
	}
}

func TestDatabaseManagerV2_MoveItem(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_move_item.json")

	persistence := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	err := persistence.SaveDataV2(models.NewDatabaseV2())
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	manager, err := createTestDatabaseManager(testDataFile)
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}

	_, err = manager.CreateFolder("source", "Source folder", "/")
	if err != nil {
		t.Fatalf("Failed to create source folder: %v", err)
	}

	_, err = manager.CreateFolder("destination", "Destination folder", "/")
	if err != nil {
		t.Fatalf("Failed to create destination folder: %v", err)
	}

	item, err := manager.CreateItem("Movable Item", "Item to move", "echo move", "/source", nil, nil)
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	err = manager.MoveItem(item.ID, "/destination")
	if err != nil {
		t.Fatalf("Failed to move item: %v", err)
	}

	movedItem, err := manager.GetItem(item.ID)
	if err != nil {
		t.Fatalf("Failed to get moved item: %v", err)
	}

	if movedItem.FolderPath != "/destination" {
		t.Errorf("Expected folder path '/destination', got %q", movedItem.FolderPath)
	}

	_, sourceItems, err := manager.GetFolderContents("/source")
	if err != nil {
		t.Fatalf("Failed to get source folder contents: %v", err)
	}

	if len(sourceItems) != 0 {
		t.Errorf("Expected 0 items in source folder, got %d", len(sourceItems))
	}

	_, destItems, err := manager.GetFolderContents("/destination")
	if err != nil {
		t.Fatalf("Failed to get destination folder contents: %v", err)
	}

	if len(destItems) != 1 {
		t.Errorf("Expected 1 item in destination folder, got %d", len(destItems))
	}
}

func TestDatabaseManagerV2_GetStatistics(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_statistics.json")

	persistence := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	err := persistence.SaveDataV2(models.NewDatabaseV2())
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	manager, err := createTestDatabaseManager(testDataFile)
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}

	_, err = manager.CreateFolder("dev", "Development", "/")
	if err != nil {
		t.Fatalf("Failed to create dev folder: %v", err)
	}

	_, err = manager.CreateFolder("scripts", "Scripts", "/dev")
	if err != nil {
		t.Fatalf("Failed to create scripts folder: %v", err)
	}

	_, err = manager.CreateItem("Item1", "First item", "echo 1", "/", []string{"tag1", "common"}, nil)
	if err != nil {
		t.Fatalf("Failed to create item1: %v", err)
	}

	_, err = manager.CreateItem("Item2", "Second item", "echo 2", "/dev", []string{"tag2", "common"}, nil)
	if err != nil {
		t.Fatalf("Failed to create item2: %v", err)
	}

	_, err = manager.CreateItem("Item3", "Third item", "echo 3", "/dev/scripts", []string{"tag1", "script"}, nil)
	if err != nil {
		t.Fatalf("Failed to create item3: %v", err)
	}

	stats := manager.GetStatistics()

	if stats["version"] != "2.0" {
		t.Errorf("Expected version '2.0', got %v", stats["version"])
	}

	if stats["total_folders"] != 2 {
		t.Errorf("Expected 2 folders, got %v", stats["total_folders"])
	}

	if stats["total_items"] != 3 {
		t.Errorf("Expected 3 items, got %v", stats["total_items"])
	}

	itemsByFolder := stats["items_by_folder"].(map[string]int)
	if itemsByFolder["/"] != 1 {
		t.Errorf("Expected 1 item in root, got %d", itemsByFolder["/"])
	}
	if itemsByFolder["/dev"] != 1 {
		t.Errorf("Expected 1 item in /dev, got %d", itemsByFolder["/dev"])
	}
	if itemsByFolder["/dev/scripts"] != 1 {
		t.Errorf("Expected 1 item in /dev/scripts, got %d", itemsByFolder["/dev/scripts"])
	}

	tagUsage := stats["tag_usage"].(map[string]int)
	if tagUsage["common"] != 2 {
		t.Errorf("Expected tag 'common' used 2 times, got %d", tagUsage["common"])
	}
	if tagUsage["tag1"] != 2 {
		t.Errorf("Expected tag 'tag1' used 2 times, got %d", tagUsage["tag1"])
	}
}

func TestDatabaseManagerV2_ValidateDatabase(t *testing.T) {
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_validate.json")

	persistence := &PersistenceService{
		dataFilePath: testDataFile,
		appName:      "test-app",
	}

	validation := NewValidationService()

	db := models.NewDatabaseV2()

	item := models.ItemV2{
		Title:       "Orphaned Item",
		Desc:        "Item in non-existent folder",
		Command:     "echo orphaned",
		FolderPath:  "/nonexistent",
		DateAdded:   time.Now(),
		DateUpdated: time.Now(),
	}
	item.GenerateID()
	db.Items = append(db.Items, item)

	err := persistence.SaveDataV2(db)
	if err != nil {
		t.Fatalf("Failed to save invalid database: %v", err)
	}

	manager, err := NewDatabaseManagerV2(persistence, validation)
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}

	issues := manager.ValidateDatabase()

	if len(issues) == 0 {
		t.Error("Expected validation issues but found none")
	}

	found := false
	for _, issue := range issues {
		if strings.Contains(issue, "Orphaned Item") && strings.Contains(issue, "non-existent folder") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected orphaned item validation issue")
	}
}
