package menu_modules

import (
	"context"

	"web100now-clients-platform/app/graph/model"
)

// Resolver handles platform video menu queries
type Resolver struct{}

// NewResolver creates a new instance of the resolver
func NewResolver() *Resolver {
	return &Resolver{}
}

// PlatformPositionsVideoMenu returns platform video menu positions grouped by tags
func (r *Resolver) PlatformPositionsVideoMenu(ctx context.Context, tagsFilter []string) ([]*model.PositionGroup, error) {
	// TODO: Implement platform video menu positions retrieval
	return []*model.PositionGroup{}, nil
}

// PlatformPositionsVideoMenuSlider returns platform video menu positions for slider
func (r *Resolver) PlatformPositionsVideoMenuSlider(ctx context.Context, tagsFilter []string, currentID string) ([]*model.Position, error) {
	// TODO: Implement platform video menu slider positions retrieval
	return []*model.Position{}, nil
}

