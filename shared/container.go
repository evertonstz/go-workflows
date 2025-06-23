package shared

import (
	"sync"
)

type Container struct {
	services map[string]interface{}
	mutex    sync.RWMutex
}

var globalContainer = &Container{
	services: make(map[string]interface{}),
}

func RegisterService(name string, service interface{}) {
	globalContainer.mutex.Lock()
	defer globalContainer.mutex.Unlock()
	globalContainer.services[name] = service
}

func GetService(name string) interface{} {
	globalContainer.mutex.RLock()
	defer globalContainer.mutex.RUnlock()
	return globalContainer.services[name]
}
