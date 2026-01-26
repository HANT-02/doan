package config

import (
	"context"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var (
	managerInstance Manager
	managerOnce     sync.Once
)

type viperManager struct {
	mu       sync.RWMutex
	watchers map[string][]chan struct{}
}

// NewManager creates a new config manager (uses singleton pattern)
func NewManager() Manager {
	managerOnce.Do(func() {
		managerInstance = &viperManager{
			watchers: make(map[string][]chan struct{}),
		}
	})
	return managerInstance
}

// GetManager returns the singleton config manager instance
func GetManager() Manager {
	if managerInstance == nil {
		return NewManager()
	}
	return managerInstance
}

func (m *viperManager) Start(ctx context.Context) {
	// Implement config watching if needed
	// For now, just a placeholder
	go func() {
		<-ctx.Done()
		// Cleanup watchers
		m.mu.Lock()
		defer m.mu.Unlock()
		for _, channels := range m.watchers {
			for _, ch := range channels {
				close(ch)
			}
		}
	}()
}

func (m *viperManager) WatchKey(key string) <-chan struct{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	ch := make(chan struct{}, 1)
	m.watchers[key] = append(m.watchers[key], ch)
	return ch
}

func (m *viperManager) UnwatchKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if channels, ok := m.watchers[key]; ok {
		for _, ch := range channels {
			close(ch)
		}
		delete(m.watchers, key)
	}
}

func (m *viperManager) Get(key string) interface{} {
	return viper.Get(key)
}

func (m *viperManager) GetString(key string) string {
	return viper.GetString(key)
}

func (m *viperManager) GetInt(key string) int {
	return viper.GetInt(key)
}

func (m *viperManager) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (m *viperManager) GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

func (m *viperManager) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func (m *viperManager) GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func (m *viperManager) UnmarshalKey(key string, rawVal interface{}) error {
	return viper.UnmarshalKey(key, rawVal)
}

func (m *viperManager) Unmarshal(rawVal interface{}) error {
	return viper.Unmarshal(rawVal)
}

func (m *viperManager) Set(key string, value interface{}) {
	viper.Set(key, value)
}

func (m *viperManager) SetDefault(key string, value interface{}) {
	viper.SetDefault(key, value)
}

func (m *viperManager) IsSet(key string) bool {
	return viper.IsSet(key)
}

func (m *viperManager) AllKeys() []string {
	return viper.AllKeys()
}

func (m *viperManager) AllSettings() map[string]interface{} {
	return viper.AllSettings()
}
