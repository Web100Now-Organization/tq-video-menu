// platform/menu_modules/slider.go
package menu_modules

import (
	"context"
	"fmt"

	"web100now-clients-platform/app/graph/model"
	"web100now-clients-platform/core/db/utils"
	"web100now-clients-platform/core/logger"
)

// GetFlatPositions повертає плоский список для слайдера
func GetFlatPositions(ctx context.Context, tagsFilter []string) ([]model.Position, error) {
	logger.LogInfo(fmt.Sprintf("[VideoMenu] GetFlatPositions called - TagsFilter: %v", tagsFilter))

	db, err := utils.GetMongoDB(ctx)
	if err != nil {
		logger.LogError("[VideoMenu] Failed to connect to MongoDB in GetFlatPositions", err)
		return nil, err
	}

	logger.LogInfo("[VideoMenu] MongoDB connection established")

	coll := db.Collection("video_menu")

	pipeline := BuildFlatPipeline(tagsFilter)
	logger.LogInfo(fmt.Sprintf("[VideoMenu] Executing flat pipeline with %d stage(s)", len(pipeline)))

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		logger.LogError("[VideoMenu] Error executing flat aggregation pipeline", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var positions []model.Position
	if err := cursor.All(ctx, &positions); err != nil {
		logger.LogError("[VideoMenu] Error decoding flat aggregation results", err)
		return nil, err
	}

	logger.LogInfo(fmt.Sprintf("[VideoMenu] GetFlatPositions completed - Found %d position(s)", len(positions)))

	return positions, nil
}
