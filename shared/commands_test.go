package shared

import (
	"errors"
	"testing"
	"time"

	"github.com/evertonstz/go-workflows/models"
	"github.com/stretchr/testify/assert"
)

func TestCopyToClipboardCmd(t *testing.T) {
	// Store original and restore after test
	originalClipboardWriteAll := ClipboardWriteAll
	defer func() { ClipboardWriteAll = originalClipboardWriteAll }()

	t.Run("success", func(t *testing.T) {
		ClipboardWriteAll = func(text string) error {
			assert.Equal(t, "test text", text, "Text passed to clipboard should be correct")
			return nil
		}

		cmd := CopyToClipboardCmd("test text")
		msg := cmd()

		assert.IsType(t, CopiedToClipboardMsg{}, msg, "Message should be CopiedToClipboardMsg on success")
	})

	t.Run("error", func(t *testing.T) {
		expectedError := errors.New("clipboard write error")
		ClipboardWriteAll = func(text string) error {
			assert.Equal(t, "error text", text, "Text passed to clipboard should be correct")
			return expectedError
		}

		cmd := CopyToClipboardCmd("error text")
		msg := cmd()

		assert.IsType(t, ErrorMsg{}, msg, "Message should be ErrorMsg on error")
		errorMsg, ok := msg.(ErrorMsg)
		assert.True(t, ok, "Failed to cast message to ErrorMsg")
		assert.Equal(t, expectedError, errorMsg.Err, "Error in message should match expected error")
	})
}

func TestSetCurrentItemCmd(t *testing.T) {
	expectedItem := models.Item{Title: "Test Item", Desc: "Description", Command: "echo test"}
	cmd := SetCurrentItemCmd(expectedItem)
	msg := cmd()

	assert.IsType(t, DidSetCurrentItemMsg{}, msg, "Message should be DidSetCurrentItemMsg")
	actualMsg, ok := msg.(DidSetCurrentItemMsg)
	assert.True(t, ok, "Failed to cast message to DidSetCurrentItemMsg")
	assert.Equal(t, expectedItem, actualMsg.Item, "Item in message should match expected item")
}

func TestDeleteCurrentItemCmd(t *testing.T) {
	expectedIndex := 5
	cmd := DeleteCurrentItemCmd(expectedIndex)
	msg := cmd()

	assert.IsType(t, DidDeleteItemMsg{}, msg, "Message should be DidDeleteItemMsg")
	actualMsg, ok := msg.(DidDeleteItemMsg)
	assert.True(t, ok, "Failed to cast message to DidDeleteItemMsg")
	assert.Equal(t, expectedIndex, actualMsg.Index, "Index in message should match expected index")
}

func TestAddNewItemCmd(t *testing.T) {
	expectedTitle := "New Command"
	expectedDesc := "This is a new command."
	expectedCommandText := "gcloud version"

	cmd := AddNewItemCmd(expectedTitle, expectedDesc, expectedCommandText)
	msg := cmd()

	assert.IsType(t, DidAddNewItemMsg{}, msg, "Message should be DidAddNewItemMsg")
	actualMsg, ok := msg.(DidAddNewItemMsg)
	assert.True(t, ok, "Failed to cast message to DidAddNewItemMsg")
	assert.Equal(t, expectedTitle, actualMsg.Title, "Title in message should match")
	assert.Equal(t, expectedDesc, actualMsg.Description, "Description in message should match")
	assert.Equal(t, expectedCommandText, actualMsg.CommandText, "CommandText in message should match")
}

func TestCloseAddNewScreenCmd(t *testing.T) {
	cmd := CloseAddNewScreenCmd()
	msg := cmd()

	assert.IsType(t, DidCloseAddNewScreenMsg{}, msg, "Message should be DidCloseAddNewScreenMsg")
}

func TestCloseConfirmationModalCmd(t *testing.T) {
	cmd := CloseConfirmationModalCmd()
	msg := cmd()

	assert.IsType(t, DidCloseConfirmationModalMsg{}, msg, "Message should be DidCloseConfirmationModalMsg")
}

func TestUpdateItemCmd(t *testing.T) {
	originalItem := models.Item{
		Title:       "Original Title",
		Desc:        "Original Description",
		Command:     "original command",
		DateAdded:   time.Now().Add(-2 * time.Hour),
		DateUpdated: time.Now().Add(-1 * time.Hour),
	}

	cmd := UpdateItemCmd(originalItem)
	timeBeforeUpdate := time.Now()
	msg := cmd()

	assert.IsType(t, DidUpdateItemMsg{}, msg, "Message should be DidUpdateItemMsg")
	actualMsg, ok := msg.(DidUpdateItemMsg)
	assert.True(t, ok, "Failed to cast message to DidUpdateItemMsg")

	updatedItem := actualMsg.Item

	// Check that original fields that shouldn't change are preserved
	assert.Equal(t, originalItem.Title, updatedItem.Title, "Title should not change")
	assert.Equal(t, originalItem.Desc, updatedItem.Desc, "Description should not change")
	assert.Equal(t, originalItem.Command, updatedItem.Command, "Command text should not change")
	assert.Equal(t, originalItem.DateAdded, updatedItem.DateAdded, "DateAdded should not change")

	// Check DateUpdated
	assert.True(t, !updatedItem.DateUpdated.IsZero(), "DateUpdated should be set (not zero)")
	assert.True(t, updatedItem.DateUpdated.After(originalItem.DateUpdated), "New DateUpdated should be after original DateUpdated")
	assert.True(t, updatedItem.DateUpdated.After(timeBeforeUpdate) || updatedItem.DateUpdated.Equal(timeBeforeUpdate), "New DateUpdated should be after or equal to the time captured before update")
}
