package random_video_main_websites

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"web100now-clients-platform/core/db/utils"
	"web100now-clients-platform/app/graph/model"
)

// Resolver підписується під Query.randomVideoMain
type Resolver struct{}

func NewResolver() *Resolver { return &Resolver{} }

// ───────────────────── публічний метод ─────────────────────

// RandomVideoMain повертає список груп за всіма тегами з plugins.config.groupOrder
func (r *Resolver) RandomVideoMain(
	ctx context.Context,
) ([]*model.RandomVideoGroup, error) {

	db, err := utils.GetMongoDB(ctx)
	if err != nil {
		return nil, err
	}

	// беремо всі теги з groupOrder
	tags, err := fetchAllTags(ctx, db)
	if err != nil {
		return nil, err
	}

	out := make([]*model.RandomVideoGroup, 0, len(tags))
	for _, tag := range tags {
		items, err := fetchAllVideosByTag(ctx, db, tag)
		if err != nil {
			return nil, err
		}
		out = append(out, &model.RandomVideoGroup{
			Tag:   tag,
			Items: items,
		})
	}
	return out, nil
}

/* ───────────────────── helpers ───────────────────── */

// fetchAllTags дістає всі ключі groupOrder у впорядкованому вигляді
func fetchAllTags(ctx context.Context, db *mongo.Database) ([]string, error) {
	var doc struct {
		Config struct {
			GroupOrder bson.M `bson:"groupOrder"` // будь-який map → bson.M
		} `bson:"config"`
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := db.Collection("plugins").
		FindOne(ctx, bson.M{"short_name": "tablq_positions_menu"}).
		Decode(&doc)
	if err != nil {
		return []string{"main-menu"}, nil // фолбек
	}

	tags := make([]string, 0, len(doc.Config.GroupOrder))
	for tag := range doc.Config.GroupOrder {
		tags = append(tags, tag)
	}
	return tags, nil
}

// fetchAllVideosByTag повертає всі документи video_menu з даним тегом у випадковому порядку
func fetchAllVideosByTag(
	ctx context.Context,
	db *mongo.Database,
	tag string,
) ([]*model.RandomPosition, error) {

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"tags", tag},
			{"videoUrlHevc", bson.D{{"$ne", nil}, {"$ne", ""}}},
			{"urlPosterPrevVideo", bson.D{{"$ne", nil}, {"$ne", ""}}},
		}}},
		{{"$sample", bson.D{{"size", 1}}}},
	}

	ctx, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()

	cur, err := db.Collection("video_menu").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var res []*model.RandomPosition
	for cur.Next(ctx) {
		var p model.RandomPosition
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		res = append(res, &p)
	}
	return res, cur.Err()
}
