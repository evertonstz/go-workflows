package main

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/shared"
)

func TestModel_Update_StateTransitions(t *testing.T) {
	tests := []struct {
		name           string
		initialState   screenState
		message        tea.Msg
		expectedState  screenState
		expectsCommand bool
	}{
		{
			name:           "Close add new screen transitions to list",
			initialState:   addNew,
			message:        shared.DidCloseAddNewScreenMsg{},
			expectedState:  newList,
			expectsCommand: false,
		},
		{
			name:           "Add new item transitions to list and persists",
			initialState:   addNew,
			message:        shared.DidAddNewItemMsg{},
			expectedState:  newList,
			expectsCommand: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal model for testing
			testModel := model{
				screenState: tt.initialState,
			}

			// Call Update
			updatedModel, cmd := testModel.Update(tt.message)
			actualModel := updatedModel.(model)

			// Check state transition
			if actualModel.screenState != tt.expectedState {
				t.Errorf("Expected state %v, got %v", tt.expectedState, actualModel.screenState)
			}

			// Check command expectation
			hasCommand := cmd != nil
			if hasCommand != tt.expectsCommand {
				t.Errorf("Expected command: %v, got command: %v", tt.expectsCommand, hasCommand)
			}
		})
	}
}

func TestModel_Update_WindowResize(t *testing.T) {
	testModel := model{
		termDimensions: termDimensions{width: 100, height: 50},
	}

	resizeMsg := tea.WindowSizeMsg{Width: 120, Height: 60}
	updatedModel, _ := testModel.Update(resizeMsg)
	actualModel := updatedModel.(model)

	if actualModel.termDimensions.width != 120 {
		t.Errorf("Expected width 120, got %d", actualModel.termDimensions.width)
	}
	if actualModel.termDimensions.height != 60 {
		t.Errorf("Expected height 60, got %d", actualModel.termDimensions.height)
	}
}

func TestModel_Update_ErrorHandling(t *testing.T) {
	testModel := model{}
	errorMsg := shared.ErrorMsg{Err: errors.New("test error")}

	_, cmd := testModel.Update(errorMsg)

	if cmd == nil {
		t.Error("Expected command for error message")
	}
}

func TestModel_IsSmallWidth(t *testing.T) {
	tests := []struct {
		width    int
		expected bool
	}{
		{width: 50, expected: true},
		{width: 100, expected: false},
		{width: 80, expected: true}, // assuming threshold is around 80
		{width: 120, expected: false},
	}

	for _, tt := range tests {
		testModel := model{
			termDimensions: termDimensions{width: tt.width},
		}

		result := testModel.isSmallWidth()
		if result != tt.expected {
			t.Errorf("Width %d: expected %v, got %v", tt.width, tt.expected, result)
		}
	}
}
