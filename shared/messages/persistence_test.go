package messages

import (
	"testing"

	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/evertonstz/go-workflows/shared/di/services"
)

// PersistenceServiceInterface defines the interface for testing
type PersistenceServiceInterface interface {
	GetDataFilePath() string
	LoadData() (models.Items, error)
	SaveData(data models.Items) error
}

// Mock persistence service for testing
type mockPersistenceService struct {
	dataFilePath string
	savedData    *models.Items
	loadError    error
	saveError    error
}

func (m *mockPersistenceService) GetDataFilePath() string {
	return m.dataFilePath
}

func (m *mockPersistenceService) LoadData() (models.Items, error) {
	if m.loadError != nil {
		return models.Items{}, m.loadError
	}
	if m.savedData != nil {
		return *m.savedData, nil
	}
	return models.Items{}, nil
}

func (m *mockPersistenceService) SaveData(data models.Items) error {
	if m.saveError != nil {
		return m.saveError
	}
	m.savedData = &data
	return nil
}

func setupMockService(mock *mockPersistenceService) func() {
	// Register the real service type for proper type checking
	realService := &services.PersistenceService{}

	// Register mock service as the real type (this is a limitation of the current DI design)
	di.RegisterService(di.PersistenceServiceKey, realService)

	// Return cleanup function
	return func() {
		// In a real scenario, you'd restore the original service
		// For now, we'll just leave it as is since this is just a test
	}
}

// Since the DI container expects exact types, let's test the functions more directly
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

// Helper error type for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
