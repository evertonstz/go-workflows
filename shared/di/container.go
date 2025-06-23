package di

import (
	"fmt"
	"sync"
)

type Container struct {
	services map[string]interface{}
	mutex    sync.RWMutex
}

var globalContainer = &Container{
	services: make(map[string]interface{}),
}

type ServiceKey int

const (
	I18nServiceKey ServiceKey = iota
	// Add other service keys here as needed
)

func RegisterService(key ServiceKey, service interface{}) {
	globalContainer.mutex.Lock()
	defer globalContainer.mutex.Unlock()
	globalContainer.services[fmt.Sprintf("%d", key)] = service
}

func GetService(key ServiceKey) interface{} {
	globalContainer.mutex.RLock()
	defer globalContainer.mutex.RUnlock()
	service, exists := globalContainer.services[fmt.Sprintf("%d", key)]
	if !exists {
		panic(fmt.Sprintf("service '%d' not found", key))
	}
	return service
}
