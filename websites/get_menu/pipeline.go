package get_menu_websites

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// BuildBasePipeline повертає набір стадій для фільтрації, додаткових полів та сортування,
// без фінального групування чи проєкції результатів.
func BuildBasePipeline(tagsFilter []string) mongo.Pipeline {
	return mongo.Pipeline{
		// Фільтруємо тільки enabled
		bson.D{{"$match", bson.D{{"status", "enabled"}}}},

		// — Приводимо _id до рядка, уніфікуємо groupName
		bson.D{{"$addFields", bson.D{
			{"id", bson.D{{"$toString", "$_id"}}},
			{"groupName", bson.D{{"$ifNull", bson.A{"$groupName", "$group_name"}}}},
		}}},
		bson.D{{"$project", bson.D{
			{"_id", 0},
			{"group_name", 0},
		}}},

		// — Lookup налаштувань груп із колекції plugins
		bson.D{{"$lookup", bson.D{
			{"from", "plugins"},
			{"pipeline", mongo.Pipeline{
				bson.D{{"$match", bson.D{{"short_name", "tablq_positions_menu"}}}},
				bson.D{{"$project", bson.D{
					{"_id", 0},
					{"groupOrder", "$config.groupOrder"},
				}}},
			}},
			{"as", "settings"},
		}}},
		bson.D{{"$unwind", "$settings"}},

		// Динамічний вибір порядку груп по першому тегу з tagsFilter
		bson.D{{"$addFields", bson.D{
			{"orderArr", bson.D{{"$let", bson.D{
				{"vars", bson.D{
					// масив пар {k,v}
					{"pairs", bson.D{{"$objectToArray", "$settings.groupOrder"}}},
					// прийшовший tagsFilter
					{"req", tagsFilter},
				}},
				{"in", bson.D{{"$let", bson.D{
					{"vars", bson.D{
						// знаходимо перший тег, що збігається з ключем об’єкта
						{"matchedKey", bson.D{{"$arrayElemAt", bson.A{
							bson.D{{"$map", bson.D{
								{"input", bson.D{{"$filter", bson.D{
									{"input", "$$pairs"},
									{"cond", bson.D{{"$in", bson.A{"$$this.k", "$$req"}}}},
								}}}},
								{"as", "p"},
								{"in", "$$p.k"},
							}}},
							0,
						}}}},
					}},
					{"in", bson.D{{
						"$getField", bson.D{
							// якщо  matchedKey==null → підмінюємо "main-menu"
							{"field", bson.D{{"$ifNull", bson.A{"$$matchedKey", "main-menu"}}}},
							{"input", "$settings.groupOrder"},
						},
					}}},
				}}}},
			}}}},
		}}},

		// — Допоміжні поля: groupIndex, hasVideo, hasGallery, tagMatchIndex
		bson.D{{"$addFields", bson.D{
			// 1) Індекс групи за plugin.config.groupOrder
			{"groupIndex", bson.D{{
				"$let", bson.D{
					{"vars", bson.D{{"orderArr", "$orderArr"}}},
					{"in", bson.D{{
						"$cond", bson.A{
							// якщо група у списку — повертаємо її індекс…
							bson.D{{"$in", bson.A{"$groupName", "$$orderArr"}}},
							bson.D{{"$indexOfArray", bson.A{"$$orderArr", "$groupName"}}},
							// …інакше розмір списку
							bson.D{{"$size", "$$orderArr"}},
						},
					}}},
				},
			}}},
			// 2) Флаг відео
			{"hasVideo", bson.D{{
				"$cond", bson.A{
					bson.D{{"$gte", bson.A{
						bson.D{{"$indexOfCP", bson.A{"$videoUrlHevc", ".mp4"}}},
						0,
					}}},
					0, 1,
				},
			}}},
			// 3) Флаг галереї
			{"hasGallery", bson.D{{
				"$cond", bson.A{
					bson.D{{"$gt", bson.A{
						bson.D{{"$size", bson.D{{"$filter", bson.D{
							{"input", bson.D{{"$ifNull", bson.A{"$galleryUrls", bson.A{}}}}},
							{"cond", bson.D{{"$and", bson.A{
								bson.D{{"$ne", bson.A{"$$this", ""}}},
								bson.D{{"$ne", bson.A{"$$this", nil}}},
							}}}},
						}}}}},
						0,
					}}},
					0, 1,
				},
			}}},
			// 4) Флаг по tagsFilter
			{"tagMatchIndex", bson.D{{
				"$let", bson.D{
					{"vars", bson.D{{"req", tagsFilter}}},
					{"in", bson.D{{
						"$cond", bson.A{
							bson.D{{"$gt", bson.A{
								bson.D{{"$size", bson.D{{"$setIntersection", bson.A{"$tags", "$$req"}}}}},
								0,
							}}},
							0, 1,
						},
					}}},
				},
			}}},
		}}},

		bson.D{{"$match", bson.D{{"tagMatchIndex", 0}}}},

		// — Сортування всередині групи
		bson.D{{"$sort", bson.D{
			{"groupIndex", 1},
			{"tagMatchIndex", 1},
			{"hasVideo", 1},
			{"hasGallery", 1},
			{"name", 1},
		}}},
	}
}

// BuildGroupPipeline додає stages для групування та фінального сортування груп
func BuildGroupPipeline(tagsFilter []string) mongo.Pipeline {
	p := BuildBasePipeline(tagsFilter)
	p = append(p,
		bson.D{{"$group", bson.D{
			{"_id", "$groupName"},
			{"items", bson.D{{"$push", "$$ROOT"}}},
			{"groupIndex", bson.D{{"$first", "$groupIndex"}}},
		}}},
		bson.D{{"$project", bson.D{
			{"groupName", "$_id"},
			{"items", 1},
			{"groupIndex", 1},
			{"_id", 0},
		}}},
		bson.D{{"$sort", bson.D{{"groupIndex", 1}}}},
	)
	return p
}

// BuildFlatPipeline повертає pipeline для плоского списку Position
func BuildFlatPipeline(tagsFilter []string) mongo.Pipeline {
	p := BuildBasePipeline(tagsFilter)
	// просто повернути відсортований slice у вигляді Position
	p = append(p,
		bson.D{{"$project", bson.D{
			{"_id", 0},
			{"id", 1},
			{"name", 1},
			{"groupName", 1},
			{"description", 1},
			{"prices", 1},
			{"unitOfMeasure", 1},
			{"usesFractionalQuantity", 1},
			{"status", 1},
			{"contentAdvisories", 1},
			{"spiceLevel", 1},
			{"nutritionInfo", 1},
			{"imageUrls", 1},
			{"videoUrlHevc", 1},
			{"urlPosterPrevVideo", 1},
			{"modifierGroups", 1},
			{"addons", 1},
			{"preparationTime", 1},
			{"availability", 1},
			{"tags", 1},
			{"dietaryLabels", 1},
			{"galleryUrls", 1},
		}}},
	)
	return p
}
