package keys

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	"github.com/stretchr/testify/assert"
)

func TestAddNewKeyMap_ShortHelp(t *testing.T) {
	expectedShortHelp := []key.Binding{AddNewKeys.Help, AddNewKeys.Close}
	actualShortHelp := AddNewKeys.ShortHelp()

	assert.Equal(t, expectedShortHelp, actualShortHelp, "AddNewKeys.ShortHelp() did not return the expected bindings")
}

func TestAddNewKeyMap_FullHelp(t *testing.T) {
	expectedFullHelp := [][]key.Binding{
		{AddNewKeys.Up, AddNewKeys.Down, AddNewKeys.Left, AddNewKeys.Right},
		{AddNewKeys.Help, AddNewKeys.Close, AddNewKeys.Submit},
	}
	actualFullHelp := AddNewKeys.FullHelp()

	assert.Equal(t, expectedFullHelp, actualFullHelp, "AddNewKeys.FullHelp() did not return the expected bindings")
}

func TestListKeyMap_ShortHelp(t *testing.T) {
	expectedShortHelp := []key.Binding{LisKeys.Help, LisKeys.Quit}
	actualShortHelp := LisKeys.ShortHelp()

	assert.Equal(t, expectedShortHelp, actualShortHelp, "LisKeys.ShortHelp() did not return the expected bindings")
}

func TestListKeyMap_FullHelp(t *testing.T) {
	// Note: LisKeys.Left and LisKeys.Right are defined in the struct but not included in FullHelp per current definition.
	// The test will reflect the current definition of FullHelp.
	expectedFullHelp := [][]key.Binding{
		{LisKeys.AddNewWorkflow, LisKeys.Delete, LisKeys.CopyWorkflow},
		{LisKeys.Up, LisKeys.Down, LisKeys.Help, LisKeys.Quit, LisKeys.Esc},
	}
	actualFullHelp := LisKeys.FullHelp()

	assert.Equal(t, expectedFullHelp, actualFullHelp, "LisKeys.FullHelp() did not return the expected bindings")
}
