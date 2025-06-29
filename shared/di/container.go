package di

import (
	"fmt"
	"sync"
)

type Container[T any] struct {
	services map[string]T
	mutex    sync.RWMutex
}

var globalContainer = &Container[any]{
	services: make(map[string]any),
}

type ServiceKey int

const (
	I18nServiceKey ServiceKey = iota
	PersistenceServiceKey
	ValidationServiceKey
	// Add other service keys here as needed
)

func RegisterService[T any](key ServiceKey, service T) {
	globalContainer.mutex.Lock()
	defer globalContainer.mutex.Unlock()
	globalContainer.services[fmt.Sprintf("%d", key)] = service
}

func GetService[T any](key ServiceKey) T {
	globalContainer.mutex.RLock()
	defer globalContainer.mutex.RUnlock()
	service, exists := globalContainer.services[fmt.Sprintf("%d", key)]
	if !exists {
		panic(fmt.Sprintf("service '%d' not found", key))
	}
	return service.(T)
}
