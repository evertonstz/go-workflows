package confirmationmodal

import (
	"testing"

	// "github.com/charmbracelet/bubbles/key" // Avoid this import for now
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// Helper to create tea.KeyMsg for testing Update method
func newKeyMsg(k string) tea.KeyMsg {
	switch k {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "left":
		return tea.KeyMsg{Type: tea.KeyLeft}
	case "right":
		return tea.KeyMsg{Type: tea.KeyRight}
	// For "h", "l", "x" and other single character keys
	default:
		if len(k) == 1 { // Treat single characters as runes
			return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
		}
		// Fallback for other strings if necessary, though not ideal for special keys
		// For this test suite, focusing on runes for single chars not explicitly handled.
		// This part might need enhancement if more complex keys are tested without key.NewMsg.
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)} // Treat as runes
	}
}

type testConfirmMsg struct{}
type testCancelMsg struct{}

var (
	dummyConfirmCmd tea.Cmd = func() tea.Msg { return testConfirmMsg{} }
	dummyCancelCmd  tea.Cmd = func() tea.Msg { return testCancelMsg{} }
)

func TestNewConfirmationModal(t *testing.T) {
	message := "Are you sure?"
	confirmButton := "Yes, proceed"
	cancelButton := "No, go back"

	m := NewConfirmationModal(message, confirmButton, cancelButton, dummyConfirmCmd, dummyCancelCmd)

	assert.NotNil(t, m, "Model should not be nil")
	assert.Equal(t, message, m.Message, "Message not set correctly")
	assert.Equal(t, confirmButton, m.ConfirmButton, "ConfirmButton not set correctly")
	assert.Equal(t, cancelButton, m.CancelButton, "CancelButton not set correctly")
	assert.NotNil(t, m.ConfirmCmd, "ConfirmCmd should not be nil")
	assert.NotNil(t, m.CancelCmd, "CancelCmd should not be nil")
	assert.Equal(t, confirm, m.selectedInput, "selectedInput should default to confirm")

	msgConfirm := m.ConfirmCmd()
	_, okConfirm := msgConfirm.(testConfirmMsg)
	assert.True(t, okConfirm, "ConfirmCmd did not return testConfirmMsg")

	msgCancel := m.CancelCmd()
	_, okCancel := msgCancel.(testCancelMsg)
	assert.True(t, okCancel, "CancelCmd did not return testCancelMsg")
}

func TestSetMessage(t *testing.T) {
	m := NewConfirmationModal("Initial", "C", "N", nil, nil)
	newMessage := "Updated message"
	m.SetMessage(newMessage)
	assert.Equal(t, newMessage, m.Message)
}

func TestSetConfirmButtonLabel(t *testing.T) {
	m := NewConfirmationModal("Msg", "Initial Confirm", "N", nil, nil)
	newLabel := "Updated Confirm"
	m.SetConfirmButtonLabel(newLabel)
	assert.Equal(t, newLabel, m.ConfirmButton)
}

func TestSetCancelButtonLabel(t *testing.T) {
	m := NewConfirmationModal("Msg", "C", "Initial Cancel", nil, nil)
	newLabel := "Updated Cancel"
	m.SetCancelButtonLabel(newLabel)
	assert.Equal(t, newLabel, m.CancelButton)
}

func TestUpdateConfirmationModal_Navigation(t *testing.T) {
	m := NewConfirmationModal("Navigate", "Confirm", "Cancel", nil, nil)

	assert.Equal(t, confirm, m.selectedInput, "Should start with confirm selected")

	// Navigate right (l)
	updatedModelL, _ := m.Update(newKeyMsg("l"))
	assert.Same(t, m, updatedModelL, "Update should return the same model instance")
	assert.Equal(t, cancel, m.selectedInput, "After l, cancel should be selected")

	// Navigate right again (l, boundary)
	m.Update(newKeyMsg("l"))
	assert.Equal(t, cancel, m.selectedInput, "After l at boundary, cancel should still be selected")

	// Navigate left (h)
	m.Update(newKeyMsg("h"))
	assert.Equal(t, confirm, m.selectedInput, "After h, confirm should be selected")

	// Navigate left again (h, boundary)
	m.Update(newKeyMsg("h"))
	assert.Equal(t, confirm, m.selectedInput, "After h at boundary, confirm should still be selected")
	
	// Navigate right (arrow key)
	updatedModelRight, _ := m.Update(newKeyMsg("right"))
	assert.Same(t, m, updatedModelRight, "Update should return the same model instance")
	assert.Equal(t, cancel, m.selectedInput, "After KeyRight, cancel should be selected")

	// Navigate left (arrow key)
	updatedModelLeft, _ := m.Update(newKeyMsg("left"))
	assert.Same(t, m, updatedModelLeft, "Update should return the same model instance")
	assert.Equal(t, confirm, m.selectedInput, "After KeyLeft, confirm should be selected")
}

func TestUpdateConfirmationModal_ActionDispatch(t *testing.T) {
	m := NewConfirmationModal("Dispatch", "OK", "Abort", dummyConfirmCmd, dummyCancelCmd)

	// Test Confirm Action
	m.selectedInput = confirm
	updatedModelConfirm, cmdConfirm := m.Update(newKeyMsg("enter"))
	assert.Same(t, m, updatedModelConfirm, "Model instance should be the same for confirm action")
	assert.NotNil(t, cmdConfirm, "Command should not be nil for confirm action")
	msgConfirm := cmdConfirm()
	_, okConfirm := msgConfirm.(testConfirmMsg)
	assert.True(t, okConfirm, "Expected testConfirmMsg for confirm action")

	// Test Cancel Action
	m.selectedInput = cancel
	updatedModelCancel, cmdCancel := m.Update(newKeyMsg("enter"))
	assert.Same(t, m, updatedModelCancel, "Model instance should be the same for cancel action")
	assert.NotNil(t, cmdCancel, "Command should not be nil for cancel action")
	msgCancel := cmdCancel()
	_, okCancel := msgCancel.(testCancelMsg)
	assert.True(t, okCancel, "Expected testCancelMsg for cancel action")
}

func TestUpdateConfirmationModal_UnhandledKey(t *testing.T) {
	m := NewConfirmationModal("Unhandled", "Yes", "No", nil, nil)
	initialSelected := m.selectedInput

	updatedModel, cmd := m.Update(newKeyMsg("x")) 

	assert.Same(t, m, updatedModel, "Model instance should be the same")
	assert.Equal(t, initialSelected, m.selectedInput, "selectedInput should not change on unhandled key")
	assert.Nil(t, cmd, "Command should be nil for unhandled key")
}

func TestInitConfirmationModal(t *testing.T) {
	m := NewConfirmationModal("Test Init", "Confirm", "Cancel", nil, nil)
	cmd := m.Init()
	assert.Nil(t, cmd, "Init should return a nil command")
}
