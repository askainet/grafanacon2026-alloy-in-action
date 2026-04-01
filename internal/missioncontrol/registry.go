package missioncontrol

import (
	"fmt"
	"sync"
)

// Registry stores all available missions
type Registry struct {
	mu       sync.RWMutex
	missions map[string]Mission
}

// NewRegistry creates a new mission registry
func NewRegistry() *Registry {
	return &Registry{
		missions: make(map[string]Mission),
	}
}

// Register adds a mission to the registry
func (r *Registry) Register(mission Mission) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.missions[mission.ID()] = mission
}

// Get retrieves a mission by ID
func (r *Registry) Get(id string) (Mission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mission, ok := r.missions[id]
	if !ok {
		return nil, fmt.Errorf("mission not found: %s", id)
	}

	return mission, nil
}

// List returns all registered missions
func (r *Registry) List() []Mission {
	r.mu.RLock()
	defer r.mu.RUnlock()

	missions := make([]Mission, 0, len(r.missions))
	for _, m := range r.missions {
		missions = append(missions, m)
	}

	return missions
}
