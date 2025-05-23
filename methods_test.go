package main

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/components/list" // For list.NewMyItem, used by test TestMyItemConstructorWorksViaListPackage
	"github.com/evertonstz/go-workflows/components/persist"
	"github.com/evertonstz/go-workflows/models"
	commandlist "github.com/evertonstz/go-workflows/screens/command_list"
	"github.com/evertonstz/go-workflows/shared" // For shared.ErrorMsg
	"github.com/stretchr/testify/assert"
)

func TestGetHelpKeys(t *testing.T) {
	t.Run("addNew state", func(t *testing.T) {
		m := model{screenState: addNew} // model is from model.go (main package)
		expectedKeys := keys.AddNewKeys
		actualKeys := m.getHelpKeys() // getHelpKeys is on main.model
		assert.Equal(t, expectedKeys, actualKeys, "Should return AddNewKeys for addNew state")
	})

	t.Run("newList state", func(t *testing.T) {
		m := model{screenState: newList}
		expectedKeys := keys.LisKeys
		actualKeys := m.getHelpKeys()
		assert.Equal(t, expectedKeys, actualKeys, "Should return LisKeys for newList state")
	})
}

func TestIsSmallWidth(t *testing.T) {
	testCases := []struct {
		name     string
		width    int
		expected bool
	}{
		{"width 80", 80, true},
		{"width 99", 99, true}, // Edge case: less than 100
		{"width 100", 100, false},
		{"width 120", 120, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use the new() constructor from model.go to get a properly initialized model
			// then override termDimensions for this specific test.
			m := new()
			m.termDimensions = termDimensions{width: tc.width, height: 50} // height is arbitrary here
			assert.Equal(t, tc.expected, m.isSmallWidth())
		})
	}
}

func TestPersistItems(t *testing.T) {
	// 1. Setup: Create model and populate listScreen
	m := new() // Use the constructor from model.go. It initializes listScreen.
	m.persistPath = "test_persist_data.json"

	// Prepare test data (models.Item for LoadedDataFileMsg)
	now := time.Now()
	testModelsItems := []models.Item{
		{Title: "Item 1", Desc: "Desc 1", Command: "cmd1", DateAdded: now.Add(-time.Hour), DateUpdated: now.Add(-30 * time.Minute)},
		{Title: "Item 2", Desc: "Desc 2", Command: "cmd2", DateAdded: now.Add(-2 * time.Hour), DateUpdated: now.Add(-time.Hour)},
	}

	// Populate listScreen using LoadedDataFileMsg
	loadMsg := persist.LoadedDataFileMsg{
		Items: models.Items{Items: testModelsItems},
	}

	// The Update method of commandlist.Model will receive this message.
	// commandlist.Model's Update delegates to its internal list.Model's Update.
	// list.Model's Update handles LoadedDataFileMsg and populates bubbles/list.Model using SetItems.
	// This relies on the fix in components/list/list.go Update method to return early.
	updatedListScreenModel, cmdLoad := m.listScreen.Update(loadMsg)
	m.listScreen = updatedListScreenModel.(commandlist.Model) // Update the model's listScreen
	if cmdLoad != nil {
		_ = cmdLoad() // Process any command returned by Update, if necessary for state
	}

	// Diagnostic assertion: Check if items are in listScreen immediately after Update
	// This assertion passed in previous debugging sessions, confirming items are loaded.
	assert.Equal(t, len(testModelsItems), len(m.listScreen.GetAllItems()), "Items should be populated in listScreen after Update")

	// 2. Mock persist.PersistDataFunc (exported var from components/persist/persist.go)
	originalPersistFunc := persist.PersistDataFunc
	var capturedPath string
	var capturedData models.Items
	persist.PersistDataFunc = func(path string, data models.Items) tea.Cmd {
		capturedPath = path
		capturedData = data
		return func() tea.Msg { return persist.PersistedFileMsg{} }
	}
	defer func() { persist.PersistDataFunc = originalPersistFunc }()

	// 3. Call persistItems (method on main.model)
	persistCmd := m.persistItems()

	// 4. Execute the command and verify message for success case
	msg := persistCmd()
	assert.IsType(t, persist.PersistedFileMsg{}, msg, "Message should be PersistedFileMsg on success")

	// 5. Verify captured arguments
	assert.Equal(t, m.persistPath, capturedPath, "Captured path should match model's persistPath")
	assert.Equal(t, len(testModelsItems), len(capturedData.Items), "Number of items in captured data should match input")

	// Compare item by item. persistItems converts list.MyItem back to models.Item.
	// So capturedData.Items should match testModelsItems.
	for i, expectedItem := range testModelsItems {
		actualItem := capturedData.Items[i]
		assert.Equal(t, expectedItem.Title, actualItem.Title, "Item %d Title mismatch", i)
		assert.Equal(t, expectedItem.Desc, actualItem.Desc, "Item %d Desc mismatch", i)
		assert.Equal(t, expectedItem.Command, actualItem.Command, "Item %d Command mismatch", i)
		assert.True(t, expectedItem.DateAdded.Equal(actualItem.DateAdded), "DateAdded should match for item %d", i)
		assert.True(t, expectedItem.DateUpdated.Equal(actualItem.DateUpdated), "DateUpdated should match for item %d", i)
	}

	// Test error case for persist.PersistDataFunc
	expectedErrorMessage := "mock persist error"
	// persist.PersistListData returns shared.ErrorMsg on error
	expectedErrorMsg := shared.ErrorMsg{Err: errors.New(expectedErrorMessage)}

	persist.PersistDataFunc = func(path string, data models.Items) tea.Cmd {
		return func() tea.Msg { return expectedErrorMsg }
	}

	persistCmdError := m.persistItems()
	errorMsg := persistCmdError() // Execute the command to get the message
	assert.IsType(t, shared.ErrorMsg{}, errorMsg, "Message should be shared.ErrorMsg on persistence error")
	// Compare error messages as direct error comparison might fail due to different instances
	assert.Equal(t, expectedErrorMsg.Err.Error(), errorMsg.(shared.ErrorMsg).Err.Error(), "Error message should match expected error")
}

// TestMyItemConstructorWorksViaListPackage confirms that list.NewMyItem is available and works.
// This is a sanity check for the test setup itself, as persistItems test relies on list.MyItem creation.
func TestMyItemConstructorWorksViaListPackage(t *testing.T) {
	title := "Test MyItem Constructor"
	desc := "MyItem Constructor Description"
	command := "echo 'Hello from MyItem Constructor'"
	now := time.Now()
	// This uses the actual list.NewMyItem from the imported package components/list
	myItem := list.NewMyItem(title, desc, command, now, now)

	assert.Equal(t, title, myItem.Title())
	assert.Equal(t, desc, myItem.Description())
	assert.Equal(t, command, myItem.Command())
	assert.True(t, now.Equal(myItem.DateAdded()))
	assert.True(t, now.Equal(myItem.DateUpdated()))
}
