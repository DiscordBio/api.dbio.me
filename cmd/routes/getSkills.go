package routes

import (
	"context"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.dbio.me/cmd/configuration"
	"go.dbio.me/cmd/middlewares"
	"go.dbio.me/cmd/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSkills(c *fiber.Ctx) error {
	DConfig := configuration.GetConfig()
	var collectionData = DConfig.Collection

	var db *mongo.Client = c.Locals("db").(*mongo.Client)
	skills := db.Database("dbio").Collection(collectionData)

	var skillsList []types.Skills

	findOptions := options.Find()
	findOptions.SetProjection(bson.M{"_id": 0, "skills": 1})

	cursor, err := skills.Find(context.Background(), bson.M{}, findOptions)

	if err != nil {
		return c.Status(500).JSON(middlewares.NewError("Internal server error"))
	}

	for cursor.Next(context.Background()) {
		var user types.Entity
		cursor.Decode(&user)

		for _, skill := range user.Skills {
			names := []string{}
			for i, s := range skillsList {
				if s.Name == skill {
					skillsList[i].Count++
					names = append(names, s.Name)

					continue
				}
			}

			if !middlewares.Contains(names, skill) {
				skillsList = append(skillsList, types.Skills{Name: skill, Count: 1, Slug: strings.ReplaceAll(strings.ToLower(skill), " ", "-")})
			}
		}
	}

	sort.Slice(skillsList, func(i, j int) bool {
		return skillsList[i].Count > skillsList[j].Count
	})

	return c.Status(200).JSON(middlewares.NewSuccess("Skills fetched", skillsList))
}
