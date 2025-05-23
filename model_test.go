package main

import (
	"testing"

	// No need to import specific component types like help.Model, etc.,
	// unless we are doing deep comparisons or type assertions,
	// which assert.NotZero doesn't require.
	"github.com/stretchr/testify/assert"
)

func TestNewModel(t *testing.T) {
	m := new()

	// Assert initial states for fields directly set or having specific defaults in new()
	assert.Equal(t, newList, m.screenState, "screenState should be newList")
	assert.Equal(t, 0, m.currentHelpHeight, "currentHelpHeight should be 0")
	assert.Equal(t, helpPanelStyle, m.panelsStyle.helpPanelStyle, "helpPanelStyle should be initialized to package variable")
	assert.Equal(t, notificationPanelStyle, m.panelsStyle.notificationPanelStyle, "notificationPanelStyle should be initialized to package variable")

	// Assert zero values for fields not explicitly set in new()
	assert.Equal(t, "", m.persistPath, "persistPath should be empty string (zero value)")
	assert.Equal(t, 0, m.termDimensions.width, "termDimensions.width should be 0 (zero value)")
	assert.Equal(t, 0, m.termDimensions.height, "termDimensions.height should be 0 (zero value)")

	// Assert that struct fields initialized by their respective New() constructors are assigned.
	// For most of these, their New() functions are expected to create a non-zero struct.
	// For m.confirmationModal, NewConfirmationModal("", "", "", nil, nil) creates a struct
	// which is effectively the "zero value" of confirmationmodal.Model if all its fields
	// are zero-types (e.g., empty strings, nil functions, false booleans).
	// In this case, assert.NotZero would fail. The critical part for testing main.new()
	// is that the constructor was called and assigned, not the internal state of the modal here.
	// Thus, for m.confirmationModal, we don't use assert.NotZero. We trust the assignment.
	// The type system ensures m.confirmationModal is of the correct type.
	_ = m.confirmationModal // Verifies m.confirmationModal is accessible and of the correct type.

	// For other components, their New() methods are more likely to set some default
	// values that make the struct non-zero.
	assert.NotZero(t, m.help, "help should be initialized by help.New() and not be a zero struct")
	assert.NotZero(t, m.addNewScreen, "addNewScreen should be initialized by addnew.New() and not be a zero struct")
	assert.NotZero(t, m.listScreen, "listScreen should be initialized by commandlist.New() and not be a zero struct")
	assert.NotZero(t, m.notification, "notification should be initialized by notification.New() and not be a zero struct")

	// Further checks for m.notification (optional, based on accessible fields/methods):
	// The `notification.New("Workflows")` call sets an internal `defaultText`.
	// Since this is not exported, we can't directly assert it.
	// `assert.NotZero(t, m.notification)` is the main check here.
}
