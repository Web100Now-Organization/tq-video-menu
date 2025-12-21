package menu_modules

import (
	"context"
	"errors"
	"fmt"

	"web100now-clients-platform/app/graph/model"
	"web100now-clients-platform/core/logger"
	"web100now-clients-platform/core/middleware/security"
)

// Resolver — конструктор плагіна
type Resolver struct{}

func NewResolver() *Resolver {
	return &Resolver{}
}

// PlatformPositionsVideoMenu — platform версія для /api/platform/v1
func (r *Resolver) PlatformPositionsVideoMenu(ctx context.Context, tagsFilter []string) ([]*model.PositionGroup, error) {
	// Перевірка OAuth аутентифікації
	if err := security.RequireOAuthAuth(ctx); err != nil {
		logger.LogError("[VideoMenu Platform] OAuth authentication required", err)
		return nil, errors.New("OAuth authentication required for platform endpoints")
	}

	// Перевірка scope (якщо потрібно)
	if err := security.RequireScope(ctx, "read"); err != nil {
		logger.LogError("[VideoMenu Platform] Required scope missing", err)
		return nil, fmt.Errorf("insufficient permissions: %w", err)
	}

	oauthUser := security.GetOAuthUserData(ctx)
	logger.LogInfo(fmt.Sprintf("[VideoMenu Platform] PlatformPositionsVideoMenu called - UserID: %s, ClientID: %s, TagsFilter: %v (%d tags)",
		oauthUser.UserID, oauthUser.ClientID, tagsFilter, len(tagsFilter)))

	groups, err := GetMenuPositions(ctx, tagsFilter)
	if err != nil {
		logger.LogError(fmt.Sprintf("[VideoMenu Platform] Error fetching menu positions - Tags: %v", tagsFilter), err)
		return nil, err
	}

	logger.LogInfo(fmt.Sprintf("[VideoMenu Platform] Successfully fetched %d position group(s)", len(groups)))

	out := make([]*model.PositionGroup, len(groups))
	for i := range groups {
		out[i] = &groups[i]
	}

	logger.LogInfo(fmt.Sprintf("[VideoMenu Platform] PlatformPositionsVideoMenu completed - Returned %d group(s)", len(out)))

	return out, nil
}

// PlatformPositionsVideoMenuSlider — platform версія для /api/platform/v1
func (r *Resolver) PlatformPositionsVideoMenuSlider(
	ctx context.Context,
	tagsFilter []string,
	currentID string,
) ([]*model.Position, error) {
	// Перевірка OAuth аутентифікації
	if err := security.RequireOAuthAuth(ctx); err != nil {
		logger.LogError("[VideoMenu Platform] OAuth authentication required", err)
		return nil, errors.New("OAuth authentication required for platform endpoints")
	}

	// Перевірка scope
	if err := security.RequireScope(ctx, "read"); err != nil {
		logger.LogError("[VideoMenu Platform] Required scope missing", err)
		return nil, fmt.Errorf("insufficient permissions: %w", err)
	}

	oauthUser := security.GetOAuthUserData(ctx)
	logger.LogInfo(fmt.Sprintf("[VideoMenu Platform] PlatformPositionsVideoMenuSlider called - UserID: %s, ClientID: %s, TagsFilter: %v, CurrentID: %s",
		oauthUser.UserID, oauthUser.ClientID, tagsFilter, currentID))

	flat, err := GetFlatPositions(ctx, tagsFilter)
	if err != nil {
		logger.LogError(fmt.Sprintf("[VideoMenu Platform] Error fetching flat positions - Tags: %v", tagsFilter), err)
		return nil, err
	}

	logger.LogInfo(fmt.Sprintf("[VideoMenu Platform] Fetched %d flat position(s)", len(flat)))

	idx := -1
	for i, pos := range flat {
		if pos.ID != nil && *pos.ID == currentID {
			idx = i
			break
		}
	}
	if idx < 0 {
		logger.LogError(fmt.Sprintf("[VideoMenu Platform] Position not found - ID: %s", currentID),
			fmt.Errorf("position not found"))
		return nil, fmt.Errorf("position with ID %s not found", currentID)
	}

	logger.LogInfo(fmt.Sprintf("[VideoMenu Platform] Found position at index %d - ID: %s", idx, currentID))

	start := idx - 2
	if start < 0 {
		start = 0
	}
	end := idx + 3
	if end > len(flat) {
		end = len(flat)
	}

	size := end - start
	slider := make([]*model.Position, 0, size)
	for i := start; i < end; i++ {
		slider = append(slider, &flat[i])
	}

	logger.LogInfo(fmt.Sprintf("[VideoMenu Platform] PlatformPositionsVideoMenuSlider completed - Returned %d position(s) (range: %d-%d)",
		len(slider), start, end-1))

	return slider, nil
}
