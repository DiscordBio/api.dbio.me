package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.dbio.me/cmd/configuration"
	"go.dbio.me/cmd/middlewares"
	"go.dbio.me/cmd/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Like(c *fiber.Ctx) error {
	DConfig := configuration.GetConfig()
	var collectionData = DConfig.Collection

	var account string
	account = middlewares.GetUserID(c)

	if account == "" {
		return c.Status(401).JSON(middlewares.NewError("Unauthorized"))
	}

	var db *mongo.Client = c.Locals("db").(*mongo.Client)
	users := db.Database("dbio").Collection(collectionData)

	var user types.Entity

	err := users.FindOne(context.Background(), bson.M{"url": c.Params("id")}).Decode(&user)

	if err != nil {
		return c.Status(404).JSON(middlewares.NewError("User not found"))
	}

	var likes []string = user.Likes

	for _, like := range likes {
		if like == account {
			return c.Status(400).JSON(middlewares.NewError("You already liked this user", fiber.Map{"isLiked": true}))
		}
	}

	likes = append(likes, account)

	_, err = users.UpdateOne(context.Background(), bson.M{"url": c.Params("id")}, bson.M{"$set": bson.M{"likes": likes}})
	if err != nil {
		return c.Status(500).JSON(middlewares.NewError("An error occured"))
	}

	return c.Status(200).JSON(middlewares.NewSuccess("Liked", fiber.Map{"isLiked": true}))
}

func Unlike(c *fiber.Ctx) error {
	DConfig := configuration.GetConfig()
	var collectionData = DConfig.Collection

	var account string
	account = middlewares.GetUserID(c)

	if account == "" {
		return c.Status(401).JSON(middlewares.NewError("Unauthorized"))
	}

	var db *mongo.Client = c.Locals("db").(*mongo.Client)
	users := db.Database("dbio").Collection(collectionData)

	var user types.Entity

	err := users.FindOne(context.Background(), bson.M{"url": c.Params("id")}).Decode(&user)

	if err != nil {
		return c.Status(404).JSON(middlewares.NewError("User not found"))
	}

	var likes []string = user.Likes

	if middlewares.Contains(likes, account) {
		likes = middlewares.Remove(likes, account)
	} else {
		return c.Status(400).JSON(middlewares.NewError("You haven't liked this user", fiber.Map{"isLiked": false}))
	}

	_, err = users.UpdateOne(context.Background(), bson.M{"url": c.Params("id")}, bson.M{"$set": bson.M{"likes": likes}})
	if err != nil {
		return c.Status(500).JSON(middlewares.NewError("An error occured"))
	}

	return c.Status(200).JSON(middlewares.NewSuccess("Unliked", fiber.Map{"isLiked": false}))
}
