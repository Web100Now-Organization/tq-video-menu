package get_menu_websites

import (
	"context"
	"fmt"

	"web100now-clients-platform/app/graph/model"
	"web100now-clients-platform/core/logger"
)

// Resolver — конструктор плагіна
type Resolver struct{}

func NewResolver() *Resolver {
	return &Resolver{}
}

// PositionsVideoMenu — проксі із GraphQL: передає tagsFilter у бізнес-логіку
func (r *Resolver) PositionsVideoMenu(ctx context.Context, tagsFilter []string) ([]*model.PositionGroup, error) {
	logger.LogInfo(fmt.Sprintf("[VideoMenu] PositionsVideoMenu called - TagsFilter: %v (%d tags)",
		tagsFilter, len(tagsFilter)))

	// 1) отримуємо всі групи згідно базової логіки
	groups, err := GetMenuPositions(ctx, tagsFilter)
	if err != nil {
		logger.LogError(fmt.Sprintf("[VideoMenu] Error fetching menu positions - Tags: %v", tagsFilter), err)
		return nil, err
	}

	logger.LogInfo(fmt.Sprintf("[VideoMenu] Successfully fetched %d position group(s)", len(groups)))

	// 2) перетворюємо на []*model.PositionGroup
	out := make([]*model.PositionGroup, len(groups))
	for i := range groups {
		out[i] = &groups[i]
	}

	logger.LogInfo(fmt.Sprintf("[VideoMenu] PositionsVideoMenu completed - Returned %d group(s)", len(out)))

	return out, nil
}

func (r *Resolver) PositionsVideoMenuSlider(
	ctx context.Context,
	tagsFilter []string,
	currentID string,
) ([]*model.Position, error) {
	logger.LogInfo(fmt.Sprintf("[VideoMenu] PositionsVideoMenuSlider called - TagsFilter: %v, CurrentID: %s",
		tagsFilter, currentID))

	// 1. Забираємо весь відсортований список
	flat, err := GetFlatPositions(ctx, tagsFilter)
	if err != nil {
		logger.LogError(fmt.Sprintf("[VideoMenu] Error fetching flat positions - Tags: %v", tagsFilter), err)
		return nil, err
	}

	logger.LogInfo(fmt.Sprintf("[VideoMenu] Fetched %d flat position(s)", len(flat)))

	// 2. Знаходимо індекс поточної позиції
	idx := -1
	for i, pos := range flat {
		if pos.ID != nil && *pos.ID == currentID {
			idx = i
			break
		}
	}
	if idx < 0 {
		logger.LogError(fmt.Sprintf("[VideoMenu] Position not found - ID: %s", currentID),
			fmt.Errorf("position not found"))
		return nil, fmt.Errorf("position with ID %s not found", currentID)
	}

	logger.LogInfo(fmt.Sprintf("[VideoMenu] Found position at index %d - ID: %s", idx, currentID))

	// 3. Розраховуємо діапазон [idx-2 .. idx+2] у межах slice
	start := idx - 2
	if start < 0 {
		start = 0
	}
	end := idx + 3 // end-exclusive: дасть індекси idx-2, idx-1, idx, idx+1, idx+2
	if end > len(flat) {
		end = len(flat)
	}

	// 4. Формуємо зріз із поточною та сусідніми позиціями
	size := end - start
	slider := make([]*model.Position, 0, size)
	for i := start; i < end; i++ {
		slider = append(slider, &flat[i])
	}

	logger.LogInfo(fmt.Sprintf("[VideoMenu] PositionsVideoMenuSlider completed - Returned %d position(s) (range: %d-%d)",
		len(slider), start, end-1))

	return slider, nil
}
