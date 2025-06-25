package messages

import (
	"testing"

	"github.com/evertonstz/go-workflows/models"
)

func TestPersistenceCommands_Direct(t *testing.T) {
	t.Run("Test message types", func(t *testing.T) {
		// Test that message types can be created
		initMsg := InitiatedPersistionMsg{DataFile: "/test/path"}
		if initMsg.DataFile != "/test/path" {
			t.Errorf("Expected DataFile '/test/path', got %q", initMsg.DataFile)
		}

		loadMsg := LoadedDataFileMsg{Items: models.Items{}}
		if len(loadMsg.Items.Items) != 0 {
			t.Errorf("Expected 0 items, got %d", len(loadMsg.Items.Items))
		}

		persistMsg := PersistedFileMsg{}
		_ = persistMsg // Just verify it can be created
	})
}
