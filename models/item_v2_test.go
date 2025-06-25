package models

import (
	"testing"
	"time"
)

func TestItemV2(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	item := ItemV2{
		Title:       "Test Workflow",
		Desc:        "Test Description",
		Command:     "echo hello",
		DateAdded:   testTime,
		DateUpdated: testTime.Add(time.Hour),
		FolderPath:  "/dev/scripts",
		Tags:        []string{"test", "script"},
		Metadata:    map[string]string{"author": "test"},
	}

	item.GenerateID()

	if item.ID == "" {
		t.Error("Expected ID to be generated")
	}

	if item.Title != "Test Workflow" {
		t.Errorf("Expected title 'Test Workflow', got %q", item.Title)
	}

	expectedPath := "/dev/scripts/Test Workflow"
	if item.GetFullPath() != expectedPath {
		t.Errorf("Expected full path %q, got %q", expectedPath, item.GetFullPath())
	}
}

func TestItemV2_GetFullPath(t *testing.T) {
	tests := []struct {
		name     string
		item     ItemV2
		expected string
	}{
		{
			name:     "root folder",
			item:     ItemV2{Title: "Test", FolderPath: "/"},
			expected: "Test",
		},
		{
			name:     "empty folder path",
			item:     ItemV2{Title: "Test", FolderPath: ""},
			expected: "Test",
		},
		{
			name:     "subfolder",
			item:     ItemV2{Title: "Test", FolderPath: "/dev/scripts"},
			expected: "/dev/scripts/Test",
		},
		{
			name:     "folder with trailing slash",
			item:     ItemV2{Title: "Test", FolderPath: "/dev/scripts/"},
			expected: "/dev/scripts/Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.GetFullPath(); got != tt.expected {
				t.Errorf("GetFullPath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestItemV2_MatchesSearch(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	item := ItemV2{
		Title:      "Test Workflow",
		Desc:       "A test description",
		Command:    "echo hello",
		DateAdded:  testTime,
		FolderPath: "/dev/scripts",
		Tags:       []string{"test", "automation"},
	}

	tests := []struct {
		name     string
		criteria SearchCriteria
		expected bool
	}{
		{
			name:     "query matches title",
			criteria: SearchCriteria{Query: "test"},
			expected: true,
		},
		{
			name:     "query matches description",
			criteria: SearchCriteria{Query: "description"},
			expected: true,
		},
		{
			name:     "query matches command",
			criteria: SearchCriteria{Query: "echo"},
			expected: true,
		},
		{
			name:     "query no match",
			criteria: SearchCriteria{Query: "nomatch"},
			expected: false,
		},
		{
			name:     "folder path matches",
			criteria: SearchCriteria{FolderPath: "/dev/scripts"},
			expected: true,
		},
		{
			name:     "folder path no match",
			criteria: SearchCriteria{FolderPath: "/other"},
			expected: false,
		},
		{
			name:     "tag matches",
			criteria: SearchCriteria{Tags: []string{"test"}},
			expected: true,
		},
		{
			name:     "tag no match",
			criteria: SearchCriteria{Tags: []string{"nomatch"}},
			expected: false,
		},
		{
			name: "date range matches",
			criteria: SearchCriteria{
				DateFrom: &time.Time{},
				DateTo:   &testTime,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := item.MatchesSearch(tt.criteria); got != tt.expected {
				t.Errorf("MatchesSearch() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFolderV2(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	folder := FolderV2{
		Name:        "Scripts",
		Description: "Development scripts",
		Path:        "/dev/scripts",
		ParentPath:  "/dev",
		DateAdded:   testTime,
		DateUpdated: testTime,
	}

	folder.GenerateID()

	if folder.ID == "" {
		t.Error("Expected ID to be generated")
	}

	if folder.IsRoot() {
		t.Error("Expected folder not to be root")
	}

	if folder.GetDepth() != 2 {
		t.Errorf("Expected depth 2, got %d", folder.GetDepth())
	}

	if !folder.IsChildOf("/dev") {
		t.Error("Expected folder to be child of /dev")
	}
}

func TestFolderV2_IsRoot(t *testing.T) {
	tests := []struct {
		name     string
		folder   FolderV2
		expected bool
	}{
		{
			name:     "root with slash",
			folder:   FolderV2{Path: "/"},
			expected: true,
		},
		{
			name:     "root empty",
			folder:   FolderV2{Path: ""},
			expected: true,
		},
		{
			name:     "not root",
			folder:   FolderV2{Path: "/dev"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.folder.IsRoot(); got != tt.expected {
				t.Errorf("IsRoot() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFolderV2_GetDepth(t *testing.T) {
	tests := []struct {
		name     string
		folder   FolderV2
		expected int
	}{
		{
			name:     "root",
			folder:   FolderV2{Path: "/"},
			expected: 0,
		},
		{
			name:     "first level",
			folder:   FolderV2{Path: "/dev"},
			expected: 1,
		},
		{
			name:     "second level",
			folder:   FolderV2{Path: "/dev/scripts"},
			expected: 2,
		},
		{
			name:     "third level",
			folder:   FolderV2{Path: "/dev/scripts/test"},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.folder.GetDepth(); got != tt.expected {
				t.Errorf("GetDepth() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDatabaseV2_AddFolder(t *testing.T) {
	db := NewDatabaseV2()

	folder := FolderV2{
		Name:        "Scripts",
		Path:        "/dev/scripts",
		ParentPath:  "/dev",
		DateAdded:   time.Now(),
		DateUpdated: time.Now(),
	}

	err := db.AddFolder(folder)
	if err != nil {
		t.Fatalf("Failed to add folder: %v", err)
	}

	if len(db.Folders) != 1 {
		t.Errorf("Expected 1 folder, got %d", len(db.Folders))
	}

	if db.Folders[0].ID == "" {
		t.Error("Expected folder ID to be generated")
	}

	// Test duplicate path
	err = db.AddFolder(folder)
	if err == nil {
		t.Error("Expected error when adding duplicate folder path")
	}

	// Test empty path
	emptyFolder := FolderV2{Name: "Empty"}
	err = db.AddFolder(emptyFolder)
	if err == nil {
		t.Error("Expected error when adding folder with empty path")
	}
}

func TestDatabaseV2_AddItem(t *testing.T) {
	db := NewDatabaseV2()

	// Add a folder first
	folder := FolderV2{
		Name:        "Scripts",
		Path:        "/dev/scripts",
		DateAdded:   time.Now(),
		DateUpdated: time.Now(),
	}
	err := db.AddFolder(folder)
	if err != nil {
		t.Fatalf("Failed to add folder: %v", err)
	}

	item := ItemV2{
		Title:       "Test Script",
		Desc:        "A test script",
		Command:     "echo test",
		FolderPath:  "/dev/scripts",
		DateAdded:   time.Now(),
		DateUpdated: time.Now(),
	}

	err = db.AddItem(item)
	if err != nil {
		t.Fatalf("Failed to add item: %v", err)
	}

	if len(db.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(db.Items))
	}

	if db.Items[0].ID == "" {
		t.Error("Expected item ID to be generated")
	}

	// Test adding item to non-existent folder
	invalidItem := ItemV2{
		Title:       "Invalid",
		FolderPath:  "/nonexistent",
		DateAdded:   time.Now(),
		DateUpdated: time.Now(),
	}

	err = db.AddItem(invalidItem)
	if err == nil {
		t.Error("Expected error when adding item to non-existent folder")
	}
}

func TestDatabaseV2_GetItemsByFolder(t *testing.T) {
	db := NewDatabaseV2()

	// Add folders
	folder1 := FolderV2{Name: "Scripts", Path: "/scripts", DateAdded: time.Now(), DateUpdated: time.Now()}
	folder2 := FolderV2{Name: "Utils", Path: "/utils", DateAdded: time.Now(), DateUpdated: time.Now()}
	if err := db.AddFolder(folder1); err != nil {
		t.Fatalf("Failed to add folder1: %v", err)
	}
	if err := db.AddFolder(folder2); err != nil {
		t.Fatalf("Failed to add folder2: %v", err)
	}

	// Add items
	item1 := ItemV2{Title: "Script1", FolderPath: "/scripts", DateAdded: time.Now(), DateUpdated: time.Now()}
	item2 := ItemV2{Title: "Script2", FolderPath: "/scripts", DateAdded: time.Now(), DateUpdated: time.Now()}
	item3 := ItemV2{Title: "Util1", FolderPath: "/utils", DateAdded: time.Now(), DateUpdated: time.Now()}

	if err := db.AddItem(item1); err != nil {
		t.Fatalf("Failed to add item1: %v", err)
	}
	if err := db.AddItem(item2); err != nil {
		t.Fatalf("Failed to add item2: %v", err)
	}
	if err := db.AddItem(item3); err != nil {
		t.Fatalf("Failed to add item3: %v", err)
	}

	scriptsItems := db.GetItemsByFolder("/scripts")
	if len(scriptsItems) != 2 {
		t.Errorf("Expected 2 items in /scripts, got %d", len(scriptsItems))
	}

	utilsItems := db.GetItemsByFolder("/utils")
	if len(utilsItems) != 1 {
		t.Errorf("Expected 1 item in /utils, got %d", len(utilsItems))
	}
}

func TestDatabaseV2_Search(t *testing.T) {
	db := NewDatabaseV2()
	testTime := time.Now()

	// Add folders
	folder := FolderV2{Name: "Scripts", Path: "/scripts", DateAdded: testTime, DateUpdated: testTime}
	if err := db.AddFolder(folder); err != nil {
		t.Fatalf("Failed to add folder: %v", err)
	}

	// Add items
	item1 := ItemV2{
		Title:       "Test Script",
		Desc:        "A test script",
		Command:     "echo test",
		FolderPath:  "/scripts",
		Tags:        []string{"test"},
		DateAdded:   testTime,
		DateUpdated: testTime,
	}
	item2 := ItemV2{
		Title:       "Deploy Script",
		Desc:        "Deploy application",
		Command:     "deploy.sh",
		FolderPath:  "/scripts",
		Tags:        []string{"deploy"},
		DateAdded:   testTime,
		DateUpdated: testTime,
	}

	if err := db.AddItem(item1); err != nil {
		t.Fatalf("Failed to add item1: %v", err)
	}
	if err := db.AddItem(item2); err != nil {
		t.Fatalf("Failed to add item2: %v", err)
	}

	// Test search by query
	result := db.Search(SearchCriteria{Query: "test"})
	if len(result.Items) != 1 {
		t.Errorf("Expected 1 item with 'test' query, got %d", len(result.Items))
	}

	// Test search by folder
	result = db.Search(SearchCriteria{FolderPath: "/scripts"})
	if len(result.Items) != 2 {
		t.Errorf("Expected 2 items in /scripts folder, got %d", len(result.Items))
	}

	// Test search by tag
	result = db.Search(SearchCriteria{Tags: []string{"deploy"}})
	if len(result.Items) != 1 {
		t.Errorf("Expected 1 item with 'deploy' tag, got %d", len(result.Items))
	}
}

func TestMigrateV1ToV2(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	v1Data := Items{
		Items: []Item{
			{
				Title:       "Old Workflow",
				Desc:        "Old Description",
				Command:     "echo old",
				DateAdded:   testTime,
				DateUpdated: testTime,
			},
		},
	}

	v2Db := MigrateV1ToV2(v1Data)

	if v2Db.Version != "2.0" {
		t.Errorf("Expected version '2.0', got %q", v2Db.Version)
	}

	if len(v2Db.Items) != 1 {
		t.Errorf("Expected 1 item after migration, got %d", len(v2Db.Items))
	}

	item := v2Db.Items[0]
	if item.Title != "Old Workflow" {
		t.Errorf("Expected title 'Old Workflow', got %q", item.Title)
	}

	if item.FolderPath != "/" {
		t.Errorf("Expected folder path '/', got %q", item.FolderPath)
	}

	if item.ID == "" {
		t.Error("Expected ID to be generated during migration")
	}
}

func TestDatabaseV2_ToV1(t *testing.T) {
	db := NewDatabaseV2()
	testTime := time.Now()

	// Add a folder first
	folder := FolderV2{
		Name:        "Scripts",
		Path:        "/scripts",
		DateAdded:   testTime,
		DateUpdated: testTime,
	}
	err := db.AddFolder(folder)
	if err != nil {
		t.Fatalf("Failed to add folder: %v", err)
	}

	item := ItemV2{
		Title:       "New Workflow",
		Desc:        "New Description",
		Command:     "echo new",
		FolderPath:  "/scripts",
		DateAdded:   testTime,
		DateUpdated: testTime,
	}

	err = db.AddItem(item)
	if err != nil {
		t.Fatalf("Failed to add item: %v", err)
	}

	v1Data := db.ToV1()

	if len(v1Data.Items) != 1 {
		t.Errorf("Expected 1 item after conversion, got %d", len(v1Data.Items))
		return
	}

	v1Item := v1Data.Items[0]
	if v1Item.Title != "New Workflow" {
		t.Errorf("Expected title 'New Workflow', got %q", v1Item.Title)
	}
}

func TestDatabaseV2_UpdateAndDelete(t *testing.T) {
	db := NewDatabaseV2()
	testTime := time.Now()

	// Add folder
	folder := FolderV2{Name: "Scripts", Path: "/scripts", DateAdded: testTime, DateUpdated: testTime}
	if err := db.AddFolder(folder); err != nil {
		t.Fatalf("Failed to add folder: %v", err)
	}

	// Add item
	item := ItemV2{
		Title:       "Original Title",
		Desc:        "Original Description",
		Command:     "echo original",
		FolderPath:  "/scripts",
		DateAdded:   testTime,
		DateUpdated: testTime,
	}
	if err := db.AddItem(item); err != nil {
		t.Fatalf("Failed to add item: %v", err)
	}

	// Get the generated ID
	itemID := db.Items[0].ID

	// Test update
	updatedItem := item
	updatedItem.Title = "Updated Title"
	err := db.UpdateItem(itemID, updatedItem)
	if err != nil {
		t.Fatalf("Failed to update item: %v", err)
	}

	retrieved, found := db.GetItemByID(itemID)
	if !found {
		t.Fatal("Item not found after update")
	}

	if retrieved.Title != "Updated Title" {
		t.Errorf("Expected updated title 'Updated Title', got %q", retrieved.Title)
	}

	// Test delete
	err = db.DeleteItem(itemID)
	if err != nil {
		t.Fatalf("Failed to delete item: %v", err)
	}

	if len(db.Items) != 0 {
		t.Errorf("Expected 0 items after deletion, got %d", len(db.Items))
	}

	// Test delete folder with items (should fail)
	item2 := ItemV2{Title: "Test", FolderPath: "/scripts", DateAdded: testTime, DateUpdated: testTime}
	if err := db.AddItem(item2); err != nil {
		t.Fatalf("Failed to add item2: %v", err)
	}

	err = db.DeleteFolder("/scripts")
	if err == nil {
		t.Error("Expected error when deleting folder with items")
	}

	// Delete item first, then folder
	if err := db.DeleteItem(db.Items[0].ID); err != nil {
		t.Fatalf("Failed to delete item: %v", err)
	}
	err = db.DeleteFolder("/scripts")
	if err != nil {
		t.Fatalf("Failed to delete empty folder: %v", err)
	}
}
