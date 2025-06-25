package di

import (
	"testing"
)

// Mock service for testing
type mockService struct {
	name string
}

func TestRegisterAndGetService(t *testing.T) {
	// Create a new container for testing to avoid conflicts
	testContainer := &Container[any]{
		services: make(map[string]any),
	}

	// Temporarily replace global container
	originalContainer := globalContainer
	globalContainer = testContainer
	defer func() {
		globalContainer = originalContainer
	}()

	// Test service
	testService := &mockService{name: "test"}

	// Register service
	RegisterService(I18nServiceKey, testService)

	// Get service
	retrievedService := GetService[*mockService](I18nServiceKey)

	if retrievedService == nil {
		t.Fatal("Expected service to be retrieved, got nil")
	}

	if retrievedService.name != "test" {
		t.Errorf("Expected service name 'test', got %q", retrievedService.name)
	}
}

func TestGetService_NotFound(t *testing.T) {
	// Create a new container for testing
	testContainer := &Container[any]{
		services: make(map[string]any),
	}

	// Temporarily replace global container
	originalContainer := globalContainer
	globalContainer = testContainer
	defer func() {
		globalContainer = originalContainer
	}()

	// Test panic recovery
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when getting non-existent service")
		}
	}()

	// Try to get non-existent service (should panic)
	GetService[*mockService](PersistenceServiceKey)
}

func TestRegisterService_Overwrite(t *testing.T) {
	// Create a new container for testing
	testContainer := &Container[any]{
		services: make(map[string]any),
	}

	// Temporarily replace global container
	originalContainer := globalContainer
	globalContainer = testContainer
	defer func() {
		globalContainer = originalContainer
	}()

	// Register first service
	firstService := &mockService{name: "first"}
	RegisterService(I18nServiceKey, firstService)

	// Register second service (should overwrite)
	secondService := &mockService{name: "second"}
	RegisterService(I18nServiceKey, secondService)

	// Get service
	retrievedService := GetService[*mockService](I18nServiceKey)

	if retrievedService.name != "second" {
		t.Errorf("Expected service name 'second', got %q", retrievedService.name)
	}
}

func TestServiceKeys(t *testing.T) {
	// Test that service keys are different
	if I18nServiceKey == PersistenceServiceKey {
		t.Error("Service keys should be different")
	}
}
