package state

import (
	"os"
	"sync"
)

// Manager provides maintenance and lifecycle handling for AppState changes.
type Manager struct {
	state    *AppState
	dirty    bool
	savePath string
	mutex    sync.Mutex
}

// NewStateManager returns an instance of Manager that handles an instance of AppState.
func NewStateManager(state *AppState, savePath string) *Manager {
	m := new(Manager)
	m.state = state
	m.dirty = false
	m.savePath = savePath

	return m
}

// Get returns the AppState handled by this Manager.
func (m *Manager) Get() *AppState {
	return m.state
}

// SetDirty flags that the current app state has changed and should be saved to disk.
func (m *Manager) SetDirty() {
	m.dirty = true
}

// Shutdown flushes any pending state updates to disk.
func (m *Manager) Shutdown() error {
	if !m.dirty {
		return nil
	}

	data, err := m.state.Serialize()
	if err != nil {
		return err
	}

	return os.WriteFile(m.savePath, data, 0644)
}
