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

func TestApp_FullOutput(t *testing.T) {
	setupTestServices(t)

	updateGolden := false
	for _, arg := range os.Args {
		if arg == "-update" {
			updateGolden = true
			break
		}
	}

	m := new()

	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(120, 40),
	)

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return len(bts) > 0
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))

	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlC,
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))

	outReader := tm.FinalOutput(t)
	outBytes, err := io.ReadAll(outReader)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	goldenFile := "testdata/TestApp_FullOutput.golden"

	if updateGolden {
		err := os.WriteFile(goldenFile, outBytes, 0o644)
		if err != nil {
			t.Fatalf("Failed to update golden file: %v", err)
		}
		t.Logf("Updated golden file: %s", goldenFile)
	} else {
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

func TestApp(t *testing.T) {
	setupTestServices(t)

	m := new()

	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(120, 40),
	)

	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	tm.Send(tea.KeyMsg{Type: tea.KeyUp})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlC,
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))
}
