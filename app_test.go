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
	_ = os.Setenv("TERM", "xterm")
	_ = os.Setenv("NO_COLOR", "1")
	_ = os.Unsetenv("COLORTERM")
	lipgloss.SetColorProfile(termenv.Ascii)
}

func setupTestServices(testName string) error {
	i18nService, err := services.NewI18nServiceWithAutoDetection("locales")
	if err != nil {
		return err
	}
	di.RegisterService(di.I18nServiceKey, i18nService)

	persistenceService, err := services.NewPersistenceService("test-app-" + testName)
	if err != nil {
		return err
	}
	di.RegisterService(di.PersistenceServiceKey, persistenceService)

	return nil
}

func TestApp_InitialScreen(t *testing.T) {
	if err := setupTestServices("initial"); err != nil {
		t.Fatalf("Failed to setup test services: %v", err)
	}

	m := new()

	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(120, 40),
	)

	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("help")) ||
				bytes.Contains(bts, []byte("Add")) ||
				bytes.Contains(bts, []byte("Enter"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*2),
	)

	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlC,
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))
}

func TestApp_FullOutput(t *testing.T) {
	if err := setupTestServices("fulloutput"); err != nil {
		t.Fatalf("Failed to setup test services: %v", err)
	}

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

	lipgloss.SetColorProfile(termenv.Ascii)

	// Create initial model
	m := new()

	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(120, 40),
	)

	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlC,
	})

	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}

	teatest.RequireEqualOutput(t, out)
}

func TestApp_FinalModel(t *testing.T) {
	if err := setupTestServices("finalmodel"); err != nil {
		t.Fatalf("Failed to setup test services: %v", err)
	}

	m := new()

	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(120, 40),
	)

	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlC,
	})

	fm := tm.FinalModel(t)
	finalModel, ok := fm.(model)
	if !ok {
		t.Fatalf("Final model has wrong type: %T", fm)
	}

	if finalModel.termDimensions.width != 120 {
		t.Errorf("Expected width 120, got %d", finalModel.termDimensions.width)
	}
	if finalModel.termDimensions.height != 40 {
		t.Errorf("Expected height 40, got %d", finalModel.termDimensions.height)
	}
}

func TestApp_WindowResize(t *testing.T) {
	if err := setupTestServices("resize"); err != nil {
		t.Fatalf("Failed to setup test services: %v", err)
	}

	m := new()

	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(80, 30),
	)

	time.Sleep(100 * time.Millisecond)

	tm.Send(tea.WindowSizeMsg{
		Width:  150,
		Height: 50,
	})

	time.Sleep(100 * time.Millisecond)

	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlC,
	})

	fm := tm.FinalModel(t)
	finalModel, ok := fm.(model)
	if !ok {
		t.Fatalf("Final model has wrong type: %T", fm)
	}

	if finalModel.termDimensions.width != 150 {
		t.Errorf("Expected width 150 after resize, got %d", finalModel.termDimensions.width)
	}
	if finalModel.termDimensions.height != 50 {
		t.Errorf("Expected height 50 after resize, got %d", finalModel.termDimensions.height)
	}
}
