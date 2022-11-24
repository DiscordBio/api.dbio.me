package routes

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.dbio.me/cmd/configuration"
	"go.dbio.me/cmd/middlewares"
	"go.dbio.me/cmd/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUser(id string, c *fiber.Ctx) error {
	DConfig := configuration.GetConfig()
	var collectionData = DConfig.Collection

	var account string
	account = middlewares.GetUserID(c)

	var db *mongo.Client = c.Locals("db").(*mongo.Client)
	users := db.Database("dbio").Collection(collectionData)

	var user types.Entity

	err := users.FindOne(context.Background(), bson.M{"url": id}).Decode(&user)

	fmt.Printf("%+v", err)

	if err != nil {
		return c.Status(404).JSON(middlewares.NewError("User not found"))
	}

	if account != "" {
		if account == user.Discord.Id {
			user.IsSelf = true
		}

		user.IsLiked = middlewares.Contains(user.Likes, account)
	}

	data := fiber.Map{
		"url":        user.URL,
		"id":         user.ID,
		"discord":    user.Discord,
		"about":      user.About,
		"socials":    user.Socials,
		"occupation": user.Occupation,
		"skills":     user.Skills,
		"likes":      middlewares.Count(user.Likes),
		"createdAt":  user.CreatedAt,
		"updatedAt":  user.UpdatedAt,
		"isLiked":    user.IsLiked,
		"isSelf":     user.IsSelf,
		"views":      middlewares.Count(user.Views),
		"isPremium":  user.Premium,
		"isVerified": user.Verified,
		"language":   user.Language,
		"roles":      user.Roles,
	}

	_, err = users.UpdateOne(context.Background(), bson.M{"url": id}, bson.M{"$set": bson.M{"views": user.Views}})

	var detailedAccount types.User
	detailedAccount = middlewares.GetDetailedUser(c)

	if user.Banner != "" {
		data["banner"] = DConfig.APIUrl + "/banners/" + user.URL
	}
	if user.Avatar != "" {
		data["avatar"] = DConfig.APIUrl + "/avatars/" + user.URL
	}

	if detailedAccount.ID != "" && detailedAccount.ID == user.Discord.Id {
		data["privacy"] = user.Privacy
		data["email"] = user.Email
		data["gender"] = user.Gender
		data["birthday"] = user.Birthday
		data["location"] = user.Location

		return c.Status(200).JSON(middlewares.NewSuccess("User fetched", data))
	} else {
		if user.Privacy.IsShow {

			if user.Privacy.IsEmailPrivate == false {
				data["email"] = user.Email
			}

			if user.Privacy.IsGenderPrivate == false {
				data["gender"] = user.Gender
				data["pronouns"] = user.Pronouns
			}

			if user.Privacy.IsBirthdayPrivate == false {
				data["birthday"] = user.Birthday
			}

			if user.Privacy.IsLocationPrivate == false {
				data["location"] = user.Location
			}

			return c.JSON(middlewares.NewSuccess("User found", data))
		} else {
			return c.Status(404).JSON(middlewares.NewError("User not found"))
		}
	}
}
