package missioncontrol

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/grafana/alloy-mission-control/internal/telemetry"
)

// Controller manages mission state in a thread-safe manner
type Controller struct {
	mu              sync.RWMutex
	activeMissions  map[string]Mission
	registry        *Registry
	metrics         *telemetry.Metrics
}

// NewController creates a new mission controller
func NewController(metrics *telemetry.Metrics) *Controller {
	return &Controller{
		activeMissions: make(map[string]Mission),
		registry:       NewRegistry(),
		metrics:        metrics,
	}
}

// Registry returns the mission registry
func (c *Controller) Registry() *Registry {
	return c.registry
}

// Activate activates a mission by ID
func (c *Controller) Activate(ctx context.Context, missionID string) error {
	// Get mission from registry
	mission, err := c.registry.Get(missionID)
	if err != nil {
		return err
	}

	// Check if already active
	c.mu.RLock()
	if _, active := c.activeMissions[missionID]; active {
		c.mu.RUnlock()
		return fmt.Errorf("mission already active: %s", missionID)
	}
	c.mu.RUnlock()

	// Activate mission
	if err := mission.OnActivate(ctx); err != nil {
		return fmt.Errorf("failed to activate mission: %w", err)
	}

	// Add to active missions
	c.mu.Lock()
	c.activeMissions[missionID] = mission
	c.mu.Unlock()

	// Update metrics
	c.metrics.MissionActive.WithLabelValues(missionID).Set(1)

	slog.Info("mission activated", "component", "missioncontrol", "mission_id", missionID, "name", mission.Name())
	return nil
}

// Deactivate deactivates a mission by ID
func (c *Controller) Deactivate(ctx context.Context, missionID string) error {
	c.mu.Lock()
	mission, active := c.activeMissions[missionID]
	if !active {
		c.mu.Unlock()
		return fmt.Errorf("mission not active: %s", missionID)
	}

	delete(c.activeMissions, missionID)
	c.mu.Unlock()

	// Deactivate mission
	if err := mission.OnDeactivate(ctx); err != nil {
		slog.Error("failed to deactivate mission", "component", "missioncontrol", "mission_id", missionID, "error", err)
		// Continue anyway to ensure mission is removed from active list
	}

	// Update metrics
	c.metrics.MissionActive.WithLabelValues(missionID).Set(0)

	slog.Info("mission deactivated", "component", "missioncontrol", "mission_id", missionID, "name", mission.Name())
	return nil
}

// DeactivateAll deactivates all active missions
func (c *Controller) DeactivateAll(ctx context.Context) {
	c.mu.RLock()
	missionIDs := make([]string, 0, len(c.activeMissions))
	for id := range c.activeMissions {
		missionIDs = append(missionIDs, id)
	}
	c.mu.RUnlock()

	for _, id := range missionIDs {
		if err := c.Deactivate(ctx, id); err != nil {
			slog.Error("failed to deactivate mission during reset", "component", "missioncontrol", "mission_id", id, "error", err)
		}
	}

	slog.Info("all missions deactivated", "component", "missioncontrol")
}

// GetActiveMissions returns a list of currently active missions
func (c *Controller) GetActiveMissions() []Mission {
	c.mu.RLock()
	defer c.mu.RUnlock()

	missions := make([]Mission, 0, len(c.activeMissions))
	for _, m := range c.activeMissions {
		missions = append(missions, m)
	}

	return missions
}

// IsActive checks if a mission is currently active
func (c *Controller) IsActive(missionID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, active := c.activeMissions[missionID]
	return active
}
