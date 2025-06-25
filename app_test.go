package main

import (
	"io"
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"

	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/evertonstz/go-workflows/shared/di/services"
)

// setupTestServices initializes the DI services needed for testing
func setupTestServices(t *testing.T) {
	localesDir := "locales"
	i18nService, err := services.NewI18nServiceWithAutoDetection(localesDir)
	if err != nil {
		t.Fatalf("Error initializing i18n service: %v", err)
	}

	di.RegisterService(di.I18nServiceKey, i18nService)

	appName := "go-workflows-test"
	persistenceService, err := services.NewPersistenceService(appName)
	if err != nil {
		t.Fatalf("Error initializing persistence service: %v", err)
	}
	di.RegisterService(di.PersistenceServiceKey, persistenceService)
}

// TestApp_FullOutput tests the full application output to ensure UI stability
func TestApp_FullOutput(t *testing.T) {
	// Setup test services
	setupTestServices(t)

	// Check if we should update golden files from test args
	updateGolden := false
	for _, arg := range os.Args {
		if arg == "-update" {
			updateGolden = true
			break
		}
	}

	// Create a test model with empty state (fresh database)
	m := new()

	// Run the test using teatest
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(120, 40),
	)

	// Wait for the initial render to complete
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return len(bts) > 0
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))

	// Send a quit message to gracefully exit
	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlC,
	})

	// Wait for the program to finish
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))

	// Get the output
	outReader := tm.FinalOutput(t)
	outBytes, err := io.ReadAll(outReader)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	// Handle golden file testing
	goldenFile := "testdata/TestApp_FullOutput.golden"

	if updateGolden {
		// Update the golden file
		err := os.WriteFile(goldenFile, outBytes, 0o644)
		if err != nil {
			t.Fatalf("Failed to update golden file: %v", err)
		}
		t.Logf("Updated golden file: %s", goldenFile)
	} else {
		// Compare against the golden file
		expected, err := os.ReadFile(goldenFile)
		if err != nil {
			if os.IsNotExist(err) {
				t.Fatalf("Golden file does not exist: %s. Run with -update to create it.", goldenFile)
			}
			t.Fatalf("Failed to read golden file: %v", err)
		}

		if string(outBytes) != string(expected) {
			t.Errorf("Output does not match golden file. Run with -update to update it.\nExpected:\n%s\nGot:\n%s", string(expected), string(outBytes))
		}
	}
}

// TestApp tests basic app functionality without golden file comparison
func TestApp(t *testing.T) {
	// Setup test services
	setupTestServices(t)

	// Create a test model
	m := new()

	// Run the test using teatest
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(120, 40),
	)

	// Send some basic navigation commands
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	tm.Send(tea.KeyMsg{Type: tea.KeyUp})

	// Send a quit message to gracefully exit
	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlC,
	})

	// Wait for the program to finish
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))

	// Just verify that the test runs without panicking
	// No golden file comparison for this test
}
