package di

import (
	"testing"
)

type mockService struct {
	name string
}

func TestRegisterAndGetService(t *testing.T) {
	testContainer := &Container[any]{
		services: make(map[string]any),
	}

	originalContainer := globalContainer
	globalContainer = testContainer
	defer func() {
		globalContainer = originalContainer
	}()

	testService := &mockService{name: "test"}

	RegisterService(I18nServiceKey, testService)

	retrievedService := GetService[*mockService](I18nServiceKey)

	if retrievedService == nil {
		t.Fatal("Expected service to be retrieved, got nil")
	}

	if retrievedService.name != "test" {
		t.Errorf("Expected service name 'test', got %q", retrievedService.name)
	}
}

func TestGetService_NotFound(t *testing.T) {
	testContainer := &Container[any]{
		services: make(map[string]any),
	}

	originalContainer := globalContainer
	globalContainer = testContainer
	defer func() {
		globalContainer = originalContainer
	}()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when getting non-existent service")
		}
	}()

	GetService[*mockService](PersistenceServiceKey)
}

func TestRegisterService_Overwrite(t *testing.T) {
	testContainer := &Container[any]{
		services: make(map[string]any),
	}

	originalContainer := globalContainer
	globalContainer = testContainer
	defer func() {
		globalContainer = originalContainer
	}()

	firstService := &mockService{name: "first"}
	RegisterService(I18nServiceKey, firstService)

	secondService := &mockService{name: "second"}
	RegisterService(I18nServiceKey, secondService)

	retrievedService := GetService[*mockService](I18nServiceKey)

	if retrievedService.name != "second" {
		t.Errorf("Expected service name 'second', got %q", retrievedService.name)
	}
}

func TestServiceKeys(t *testing.T) {
	if I18nServiceKey == PersistenceServiceKey {
		t.Error("Service keys should be different")
	}
}
