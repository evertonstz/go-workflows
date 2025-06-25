package models

import (
	"testing"
	"time"
)

func TestItem(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	item := Item{
		Title:       "Test Workflow",
		Desc:        "Test Description",
		Command:     "echo hello",
		DateAdded:   testTime,
		DateUpdated: testTime.Add(time.Hour),
	}

	// Test all fields are set correctly
	if item.Title != "Test Workflow" {
		t.Errorf("Expected title 'Test Workflow', got %q", item.Title)
	}

	if item.Desc != "Test Description" {
		t.Errorf("Expected desc 'Test Description', got %q", item.Desc)
	}

	if item.Command != "echo hello" {
		t.Errorf("Expected command 'echo hello', got %q", item.Command)
	}

	if !item.DateAdded.Equal(testTime) {
		t.Errorf("Expected date added %v, got %v", testTime, item.DateAdded)
	}

	if !item.DateUpdated.Equal(testTime.Add(time.Hour)) {
		t.Errorf("Expected date updated %v, got %v", testTime.Add(time.Hour), item.DateUpdated)
	}
}

func TestItems(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	items := Items{
		Items: []Item{
			{
				Title:       "Workflow 1",
				Desc:        "Description 1",
				Command:     "echo 1",
				DateAdded:   testTime,
				DateUpdated: testTime,
			},
			{
				Title:       "Workflow 2",
				Desc:        "Description 2",
				Command:     "echo 2",
				DateAdded:   testTime.Add(time.Hour),
				DateUpdated: testTime.Add(time.Hour),
			},
		},
	}

	// Test that Items slice contains correct number of items
	if len(items.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items.Items))
	}

	// Test first item
	firstItem := items.Items[0]
	if firstItem.Title != "Workflow 1" {
		t.Errorf("Expected first item title 'Workflow 1', got %q", firstItem.Title)
	}

	// Test second item
	secondItem := items.Items[1]
	if secondItem.Title != "Workflow 2" {
		t.Errorf("Expected second item title 'Workflow 2', got %q", secondItem.Title)
	}
}

func TestEmptyItems(t *testing.T) {
	items := Items{}

	// Test empty Items
	if len(items.Items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(items.Items))
	}

	// Test nil Items slice
	if items.Items != nil {
		t.Error("Expected Items slice to be nil initially")
	}
}
