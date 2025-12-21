// platform/menu_modules/positions.go
package menu_modules

import (
	"context"
	"fmt"

	"web100now-clients-platform/app/graph/model"
	"web100now-clients-platform/core/db/utils"
	"web100now-clients-platform/core/logger"
)

// GetMenuPositions повертає згруповані позиції
func GetMenuPositions(ctx context.Context, tagsFilter []string) ([]model.PositionGroup, error) {
	logger.LogInfo(fmt.Sprintf("[VideoMenu] GetMenuPositions called - TagsFilter: %v", tagsFilter))

	db, err := utils.GetMongoDB(ctx)
	if err != nil {
		logger.LogError("[VideoMenu] Failed to connect to MongoDB in GetMenuPositions", err)
		return nil, err
	}

	logger.LogInfo("[VideoMenu] MongoDB connection established")

	coll := db.Collection("video_menu")

	pipeline := BuildGroupPipeline(tagsFilter)
	logger.LogInfo(fmt.Sprintf("[VideoMenu] Executing aggregation pipeline with %d stage(s)", len(pipeline)))

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		logger.LogError("[VideoMenu] Error executing aggregation pipeline", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []model.PositionGroup
	if err := cursor.All(ctx, &groups); err != nil {
		logger.LogError("[VideoMenu] Error decoding aggregation results", err)
		return nil, err
	}

	logger.LogInfo(fmt.Sprintf("[VideoMenu] GetMenuPositions completed - Found %d position group(s)", len(groups)))

	return groups, nil
}
