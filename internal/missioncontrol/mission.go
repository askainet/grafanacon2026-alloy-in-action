package missioncontrol

import (
	"context"
)

// Mission represents a workshop scenario that modifies telemetry generation
type Mission interface {
	// Metadata
	ID() string
	Name() string
	Description() string

	// Lifecycle hooks
	OnActivate(ctx context.Context) error
	OnDeactivate(ctx context.Context) error
}
