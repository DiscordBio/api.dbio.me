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

func GetRoles(c *fiber.Ctx) error {
	DConfig := configuration.GetConfig()
	var collectionData = DConfig.Collection
	
	var db *mongo.Client = c.Locals("db").(*mongo.Client)
	roles := db.Database("dbio").Collection(collectionData)

	var rolesList []types.Roles

	findOptions := options.Find()
	findOptions.SetProjection(bson.M{"_id": 0, "roles": 1})

	cursor, err := roles.Find(context.Background(), bson.M{}, findOptions)

	if err != nil {
		return c.Status(500).JSON(middlewares.NewError("Internal server error"))
	}

	for cursor.Next(context.Background()) {
		var user types.Entity
		cursor.Decode(&user)

		for _, role := range user.Roles {
			names := []string{}
			for i, r := range rolesList {
				if r.Name == role {
					rolesList[i].Count++
					names = append(names, r.Name)

					continue
				}
			}

			if !middlewares.Contains(names, role) {
				rolesList = append(rolesList, types.Roles{Name: role, Count: 1, Slug: strings.ReplaceAll(strings.ToLower(role), " ", "-")})
			}
		}
	}

	sort.Slice(rolesList, func(i, j int) bool {
		return rolesList[i].Count > rolesList[j].Count
	})

	return c.Status(200).JSON(middlewares.NewSuccess("Roles fetched", rolesList))
}
