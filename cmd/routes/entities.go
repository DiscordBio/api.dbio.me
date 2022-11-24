package routes

import (
	"context"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.dbio.me/cmd/configuration"
	"go.dbio.me/cmd/middlewares"
	"go.dbio.me/cmd/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Entities(c *fiber.Ctx) error {
	DConfig := configuration.GetConfig()
	var collectionData = DConfig.Collection

	var account string
	account = middlewares.GetUserID(c)

	var pagination types.Pagination = types.Pagination{
		Page:     "1",
		Limit:    "10",
		Sort:     "popular",
		Roles:    "",
		Skills:   "",
		Query:    "",
		Language: "",
	}

	if c.Query("page") != "" {
		pagination.Page = c.Query("page")
	}

	if c.Query("limit") != "" {
		pagination.Limit = c.Query("limit")
	}

	if c.Query("sort") != "" {
		pagination.Sort = c.Query("sort")
	}

	if c.Query("roles") != "" {
		pagination.Roles = c.Query("roles")
	}

	if c.Query("skills") != "" {
		pagination.Skills = c.Query("skills")
	}

	if c.Query("query") != "" {
		pagination.Query = c.Query("query")
	}

	if c.Query("language") != "" {
		pagination.Language = c.Query("language")
	}

	page, err := strconv.Atoi(pagination.Page)
	if err != nil {
		return c.Status(400).JSON(middlewares.NewError("Invalid page"))
	}

	limit, err := strconv.Atoi(pagination.Limit)
	if err != nil {
		return c.Status(400).JSON(middlewares.NewError("Invalid limit"))
	}

	if limit > 100 {
		return c.Status(400).JSON(middlewares.NewError("Limit cannot be greater than 100"))
	} else if limit < 1 {
		return c.Status(400).JSON(middlewares.NewError("Limit cannot be less than 1"))
	}

	var db *mongo.Client = c.Locals("db").(*mongo.Client)
	users := db.Database("dbio").Collection(collectionData)

	var usersList []interface{}

	findOptions := options.Find()
	findOptions.SetProjection(bson.M{"_id": 0, "id": 1, "url": 1, "discord": 1, "about": 1, "avatar": 1, "occupation": 1, "likes": 1, "language": 1, "banner": 1, "createdAt": 1, "updatedAt": 1, "views": 1, "isTeamMember": 1, "isPremium": 1, "isVerified": 1})

	filter := bson.M{
		"privacy.isShow": true,
		"$or": []bson.M{
			{"url": bson.M{"$regex": pagination.Query, "$options": "i"}},
			{"name": bson.M{"$regex": pagination.Query, "$options": "i"}},
			{"occupation": bson.M{"$regex": pagination.Query, "$options": "i"}},
			{"about": bson.M{"$regex": pagination.Query, "$options": "i"}},
			{"skills": bson.M{"$regex": pagination.Query, "$options": "i"}},
			{"roles": bson.M{"$regex": pagination.Query, "$options": "i"}},
		},
	}

	if c.Query("isTeamMembers") == "true" {
		filter["isTeamMember"] = true
	}

	if pagination.Roles != "" {
		roles := strings.Split(pagination.Roles, ",")
		filter["roles"] = bson.M{"$in": roles}
	}

	if pagination.Skills != "" {
		skills := strings.Split(pagination.Skills, ",")
		filter["skills"] = bson.M{"$in": skills}
	}

	if pagination.Language != "" {
		filter["language"] = pagination.Language
	}

	if pagination.Sort == "newest" {
		findOptions.SetSort(bson.D{{"createdAt", -1}})
	} else if pagination.Sort == "oldest" {
		findOptions.SetSort(bson.D{{"createdAt", 1}})
	} else if pagination.Sort == "popular" {
		findOptions.SetSort(bson.D{{"likes", -1}})
	} else if pagination.Sort == "trending" {
		findOptions.SetSort(bson.D{{"views", -1}})
	} else if pagination.Sort == "ascending" {
		findOptions.SetSort(bson.D{{"url", 1}})
	} else if pagination.Sort == "descending" {
		findOptions.SetSort(bson.D{{"url", -1}})
	} else {
		findOptions.SetSort(bson.D{{"likes", -1}})
	}

	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64((page - 1) * limit))

	cursor, err := users.Find(context.Background(), filter, findOptions)
	if err != nil {
		return c.Status(500).JSON(middlewares.NewError("Internal server error"))
	}

	for cursor.Next(context.Background()) {
		var user types.Entity
		cursor.Decode(&user)

		data := fiber.Map{
			"url":        user.URL,
			"id":         user.ID,
			"discord":    user.Discord,
			"about":      user.About,
			"language":   user.Language,
			"occupation": user.Occupation,
			"roles":      user.Roles,
			"likes":      middlewares.Count(user.Likes),
			"createdAt":  user.CreatedAt,
			"updatedAt":  user.UpdatedAt,
			"isLiked":    user.IsLiked,
			"isSelf":     user.IsSelf,
			"views":      middlewares.Count(user.Views),
			"isPremium":  user.Premium,
			"isVerified": user.Verified,
		}

		if user.Banner != "" {
			data["banner"] = DConfig.APIUrl + "/banners/" + user.URL
		}
		if user.Avatar != "" {
			data["avatar"] = DConfig.APIUrl + "/avatars/" + user.URL
		}

		if account != "" {
			if account == user.Discord.Id {
				data["isSelf"] = true
			}

			data["isLiked"] = middlewares.Contains(user.Likes, account)
		}

		usersList = append(usersList, data)
	}

	var count int64
	count, err = users.CountDocuments(context.Background(), filter)
	if err != nil {
		return c.Status(500).JSON(middlewares.NewError("Internal server error"))
	}

	return c.Status(200).JSON(middlewares.NewSuccess("Users fetched", fiber.Map{
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": count,
			"more":  page*limit < int(count),
		},
		"users": usersList,
	}))
}

func GetEntity(c *fiber.Ctx) error {
	return FindUser(c.Params("id"), c)
}
