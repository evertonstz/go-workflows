package main

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/muesli/termenv"

	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/evertonstz/go-workflows/shared/di/services"
)

// Force consistent color profile for CI/testing
func init() {
	// Set environment variables for consistent terminal behavior
	_ = os.Setenv("TERM", "xterm")
	_ = os.Setenv("NO_COLOR", "1")
	_ = os.Unsetenv("COLORTERM")

	// Force ASCII color profile for consistent output
	lipgloss.SetColorProfile(termenv.Ascii)
}

// setupTestServices sets up all required services for testing
func setupTestServices(testName string) error {
	// Setup I18n service
	i18nService, err := services.NewI18nServiceWithAutoDetection("locales")
	if err != nil {
		return err
	}
	di.RegisterService(di.I18nServiceKey, i18nService)

	// Setup Persistence service with unique name for each test
	persistenceService, err := services.NewPersistenceService("test-app-" + testName)
	if err != nil {
		return err
	}
	di.RegisterService(di.PersistenceServiceKey, persistenceService)

	return nil
}

// TestApp_InitialScreen tests the initial screen state and output
func TestApp_InitialScreen(t *testing.T) {
	// Setup all required services
	if err := setupTestServices("initial"); err != nil {
		t.Fatalf("Failed to setup test services: %v", err)
	}

	// Create initial model using the correct function name
	m := new()

	// Create test model with fixed terminal size
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(120, 40),
	)

	// Wait for initial output to contain expected elements
	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			// Check for presence of help text or UI elements
			return bytes.Contains(bts, []byte("help")) ||
				bytes.Contains(bts, []byte("Add")) ||
				bytes.Contains(bts, []byte("Enter"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*2),
	)

	// Send quit command
	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlC,
	})

	// Wait for the program to finish
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))
}

// TestApp_FullOutput tests the complete output of the application
func TestApp_FullOutput(t *testing.T) {
	// Setup all required services
	if err := setupTestServices("fulloutput"); err != nil {
		t.Fatalf("Failed to setup test services: %v", err)
	}
	// Force the most basic terminal environment for golden file testing
	// This ensures consistent output across local, CI, and act environments
	originalTerm := os.Getenv("TERM")
	originalColorTerm := os.Getenv("COLORTERM")

	_ = os.Setenv("TERM", "dumb")
	_ = os.Setenv("NO_COLOR", "1")
	_ = os.Unsetenv("COLORTERM")

	defer func() {
		_ = os.Setenv("TERM", originalTerm)
		if originalColorTerm != "" {
			_ = os.Setenv("COLORTERM", originalColorTerm)
		}
		_ = os.Unsetenv("NO_COLOR")
	}()

	// Force ASCII profile again to be sure
	lipgloss.SetColorProfile(termenv.Ascii)

	// Create initial model
	m := new()

	// Create test model with consistent settings
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(120, 40),
	)

	// Send quit immediately to get initial output
	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlC,
	})

	// Get final output
	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}

	// Assert using golden file (will create if doesn't exist)
	teatest.RequireEqualOutput(t, out)
}

// TestApp_FinalModel tests the final model state
func TestApp_FinalModel(t *testing.T) {
	// Setup all required services
	if err := setupTestServices("finalmodel"); err != nil {
		t.Fatalf("Failed to setup test services: %v", err)
	}

	// Create initial model
	m := new()

	// Create test model
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(120, 40),
	)

	// Send quit command
	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlC,
	})

	// Get final model
	fm := tm.FinalModel(t)
	finalModel, ok := fm.(model)
	if !ok {
		t.Fatalf("Final model has wrong type: %T", fm)
	}

	// Assert model state
	if finalModel.termDimensions.width != 120 {
		t.Errorf("Expected width 120, got %d", finalModel.termDimensions.width)
	}
	if finalModel.termDimensions.height != 40 {
		t.Errorf("Expected height 40, got %d", finalModel.termDimensions.height)
	}
}

// TestApp_WindowResize tests window resize handling
func TestApp_WindowResize(t *testing.T) {
	// Setup all required services
	if err := setupTestServices("resize"); err != nil {
		t.Fatalf("Failed to setup test services: %v", err)
	}

	// Create initial model
	m := new()

	// Create test model with initial size
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(80, 30),
	)

	// Wait for initial render
	time.Sleep(100 * time.Millisecond)

	// Send window resize message
	tm.Send(tea.WindowSizeMsg{
		Width:  150,
		Height: 50,
	})

	// Wait for resize to be processed
	time.Sleep(100 * time.Millisecond)

	// Send quit command
	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlC,
	})

	// Get final model and check dimensions were updated
	fm := tm.FinalModel(t)
	finalModel, ok := fm.(model)
	if !ok {
		t.Fatalf("Final model has wrong type: %T", fm)
	}

	// The model should have updated its terminal dimensions
	if finalModel.termDimensions.width != 150 {
		t.Errorf("Expected width 150 after resize, got %d", finalModel.termDimensions.width)
	}
	if finalModel.termDimensions.height != 50 {
		t.Errorf("Expected height 50 after resize, got %d", finalModel.termDimensions.height)
	}
}
