package get_menu

import (
	"context"

	"web100now-clients-platform/app/graph/model"
)

// Resolver handles video menu queries for websites
type Resolver struct{}

// NewResolver creates a new instance of the resolver
func NewResolver() *Resolver {
	return &Resolver{}
}

// PositionsVideoMenu returns video menu positions grouped by tags
func (r *Resolver) PositionsVideoMenu(ctx context.Context, tagsFilter []string) ([]*model.PositionGroup, error) {
	// TODO: Implement video menu positions retrieval
	return []*model.PositionGroup{}, nil
}

// PositionsVideoMenuSlider returns video menu positions for slider
func (r *Resolver) PositionsVideoMenuSlider(ctx context.Context, tagsFilter []string, currentID string) ([]*model.Position, error) {
	// TODO: Implement video menu slider positions retrieval
	return []*model.Position{}, nil
}

