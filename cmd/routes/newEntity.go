package routes

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.dbio.me/cmd/configuration"
	"go.dbio.me/cmd/middlewares"
	"go.dbio.me/cmd/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewEntity(c *fiber.Ctx) error {
	DConfig := configuration.GetConfig()
	var collectionData = DConfig.Collection

	var account string
	account = middlewares.GetUserID(c)

	if account == "" {
		return c.Status(401).JSON(middlewares.NewError("Unauthorized"))
	}

	var detailedAccount types.User
	detailedAccount = middlewares.GetDetailedUser(c)

	if detailedAccount.ID == "" {
		return c.Status(401).JSON(middlewares.NewError("Unauthorized"))
	}

	var db *mongo.Client = c.Locals("db").(*mongo.Client)
	users := db.Database("dbio").Collection(collectionData)

	body := make(map[string]interface{})
	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(400).JSON(middlewares.NewError("Invalid body"))
	}

	url := body["url"].(string)

	var thisUser bson.M

	err = users.FindOne(context.Background(), bson.M{"discord.id": account}).Decode(&thisUser)

	if err != mongo.ErrNoDocuments {

		if thisUser["url"] != url {
			if users.FindOne(context.Background(), bson.M{"url": url}).Err() != mongo.ErrNoDocuments {
				return c.Status(400).JSON(middlewares.NewError("URL already taken"))
			}
		}

		if body["about"] == "" || body["about"] == nil {
			body["about"] = thisUser["about"]
		}

		if body["url"] == "" || body["url"] == nil {
			body["url"] = thisUser["url"]
		}

		if body["banner"] != nil && body["banner"] != "" {
			if !strings.HasPrefix(body["banner"].(string), "data:image/") {
				body["banner"] = thisUser["banner"]
			} else {
				bannerId := middlewares.GenerateID()
				var imageName string
				imageName = "banner-" + bannerId

				if strings.HasPrefix(body["banner"].(string), "data:image/gif") {
					if thisUser["isPremium"] == false {
						return c.Status(400).JSON(middlewares.NewError("You need to be premium to upload a gif banner"))
					}
				}

				image, err := middlewares.UploadImage(imageName, body["banner"].(string))

				if err == nil {
					body["banner"] = image.Data.DisplayURL
				} else {
					body["banner"] = thisUser["banner"]
				}
			}
		} else {
			body["banner"] = thisUser["banner"]
		}

		if body["avatar"] != nil && body["avatar"] != "" {
			if !strings.HasPrefix(body["avatar"].(string), "data:image/") {
				body["avatar"] = thisUser["avatar"]
			} else {
				avatarId := middlewares.GenerateID()
				var imageName string
				imageName = "avatar-" + avatarId

				if strings.HasPrefix(body["avatar"].(string), "data:image/gif") {
					if thisUser["isPremium"] == false {
						return c.Status(400).JSON(middlewares.NewError("You need to be premium to upload a gif avatar"))
					}
				}

				image, err := middlewares.UploadImage(imageName, body["avatar"].(string))

				if err == nil {
					body["avatar"] = image.Data.DisplayURL
				} else {
					body["avatar"] = thisUser["avatar"]
				}
			}
		} else {
			body["avatar"] = thisUser["avatar"]
		}

		if body["socials"] != nil {
			var socialConfig types.Social = configuration.GetSocials()

			for i := 0; i < len(socialConfig); i++ {
				for j := 0; j < len(body["socials"].([]interface{})); j++ {
					if body["socials"].([]interface{})[j].(map[string]interface{})["name"] == socialConfig[i].Name {
						body["socials"].([]interface{})[j].(map[string]interface{})["icon"] = socialConfig[i].Icon
						if body["socials"].([]interface{})[j].(map[string]interface{})["url"] != "" {
							slicePath := socialConfig[i].URL
							slicePath = strings.Replace(slicePath, "{username}", "", -1)
							body["socials"].([]interface{})[j].(map[string]interface{})["id"] = socialConfig[i].ID
							body["socials"].([]interface{})[j].(map[string]interface{})["username"] = strings.Replace(body["socials"].([]interface{})[j].(map[string]interface{})["url"].(string), slicePath, "", -1)
							body["socials"].([]interface{})[j].(map[string]interface{})["url"] = slicePath + body["socials"].([]interface{})[j].(map[string]interface{})["username"].(string)
							body["socials"].([]interface{})[j].(map[string]interface{})["color"] = socialConfig[i].Color

							body["socials"] = body["socials"].([]interface{})

						}
					}
				}
			}
		}

		_, err = users.UpdateOne(context.Background(), bson.M{"discord.id": account}, bson.M{
			"$set": bson.M{
				"discord": bson.M{
					"id":            account,
					"username":      detailedAccount.Username,
					"discriminator": detailedAccount.Discriminator,
				},
				"avatar":     body["avatar"],
				"about":      body["about"],
				"url":        strings.ToLower(middlewares.RemoveSpecialChars(body["url"].(string))),
				"location":   body["location"],
				"occupation": body["occupation"],
				"gender":     body["gender"],
				"birthday":   body["birthday"],
				"language":   body["language"],
				"email":      body["email"],
				"roles":      body["roles"],
				"socials":    body["socials"],
				"skills":     body["skills"],
				"banner":     body["banner"],
				"pronouns":   body["pronouns"],
				"privacy": bson.M{
					"isShow":            body["privacy"] != nil && body["privacy"].(map[string]interface{})["isShow"] != nil && body["privacy"].(map[string]interface{})["isShow"] != "" && body["privacy"].(map[string]interface{})["isShow"] != "null",
					"isEmailPrivate":    body["privacy"] != nil && body["privacy"].(map[string]interface{})["isEmailPrivate"] != nil && body["privacy"].(map[string]interface{})["isEmailPrivate"] != "" && body["privacy"].(map[string]interface{})["isEmailPrivate"] != "null",
					"isGenderPrivate":   body["privacy"] != nil && body["privacy"].(map[string]interface{})["isGenderPrivate"] != nil && body["privacy"].(map[string]interface{})["isGenderPrivate"] != "" && body["privacy"].(map[string]interface{})["isGenderPrivate"] != "null",
					"isLocationPrivate": body["privacy"] != nil && body["privacy"].(map[string]interface{})["isLocationPrivate"] != nil && body["privacy"].(map[string]interface{})["isLocationPrivate"] != "" && body["privacy"].(map[string]interface{})["isLocationPrivate"] != "null",
					"isBirthdayPrivate": body["privacy"] != nil && body["privacy"].(map[string]interface{})["isBirthdayPrivate"] != nil && body["privacy"].(map[string]interface{})["isBirthdayPrivate"] != "" && body["privacy"].(map[string]interface{})["isBirthdayPrivate"] != "null",
				},
				"updatedAt": time.Now(),
			},
		})

		if err != nil {
			return c.Status(500).JSON(middlewares.NewError("Internal server error"))
		}

		return c.JSON(middlewares.NewSuccess("User updated", fiber.Map{}))

	} else {

		if body["about"] == nil || body["about"] == "" || body["url"] == nil || body["url"] == "" {
			return c.Status(400).JSON(middlewares.NewError("Missing required fields"))
		}

		if body["location"] == "" || body["location"] == nil {
			body["location"] = nil
		}

		if body["occupation"] == "" || body["occupation"] == nil {
			body["occupation"] = nil
		}

		if body["gender"] == "" || body["gender"] == nil {
			body["gender"] = nil
		}

		if body["birthday"] == "" || body["birthday"] == nil {
			body["birthday"] = nil
		}

		if body["language"] == "" || body["language"] == nil {
			body["language"] = nil
		}

		if body["email"] == "" || body["email"] == nil {
			body["email"] = nil
		}

		if body["privacy"] == nil {
			body["privacy"] = nil
		}

		if body["socials"] == nil {
			body["socials"] = nil
		}

		if body["roles"] == nil {
			body["roles"] = nil
		}

		if body["skills"] == nil {
			body["skills"] = nil
		}

		if body["banner"] != nil && body["banner"] != "" {
			if !strings.HasPrefix(body["banner"].(string), "data:image/") {
				body["banner"] = "https://cdn.dbio.me/assets/banner-without-text.png"
			} else {
				bannerId := middlewares.GenerateID()
				var imageName string
				imageName = "banner-" + bannerId

				if strings.HasPrefix(body["banner"].(string), "data:image/gif") {
					if thisUser["isPremium"] == false {
						return c.Status(400).JSON(middlewares.NewError("You need to be premium to upload a gif banner"))
					}
				}

				image, err := middlewares.UploadImage(imageName, body["banner"].(string))

				if err == nil {
					body["banner"] = image.Data.DisplayURL
				} else {
					body["banner"] = "https://cdn.dbio.me/assets/banner-without-text.png"
				}
			}
		} else {
			body["banner"] = "https://cdn.dbio.me/assets/banner-without-text.png"
		}

		if body["avatar"] != nil && body["avatar"] != "" {
			if !strings.HasPrefix(body["avatar"].(string), "data:image/") {
				body["avatar"] = "https://cdn.discordapp.com/avatars/" + account + "/" + detailedAccount.Avatar + ".png"
			} else {
				avatarId := middlewares.GenerateID()
				var imageName string
				imageName = "avatar-" + avatarId

				if strings.HasPrefix(body["avatar"].(string), "data:image/gif") {
					if thisUser["isPremium"] == false {
						return c.Status(400).JSON(middlewares.NewError("You need to be premium to upload a gif avatar"))
					}
				}

				image, err := middlewares.UploadImage(imageName, body["avatar"].(string))

				if err == nil {
					body["avatar"] = image.Data.DisplayURL
				} else {
					body["avatar"] = "https://cdn.discordapp.com/avatars/" + account + "/" + detailedAccount.Avatar + ".png"
				}
			}
		} else {
			body["avatar"] = "https://cdn.discordapp.com/avatars/" + account + "/" + detailedAccount.Avatar + ".png"
		}

		if body["pronouns"] == nil {
			body["pronouns"] = nil
		}

		if body["socials"] != nil {
			var socialConfig types.Social = configuration.GetSocials()

			for i := 0; i < len(socialConfig); i++ {
				for j := 0; j < len(body["socials"].([]interface{})); j++ {
					if body["socials"].([]interface{})[j].(map[string]interface{})["name"] == socialConfig[i].Name {
						body["socials"].([]interface{})[j].(map[string]interface{})["icon"] = socialConfig[i].Icon
						if body["socials"].([]interface{})[j].(map[string]interface{})["url"] != "" {
							slicePath := socialConfig[i].URL
							slicePath = strings.Replace(slicePath, "{username}", "", -1)
							body["socials"].([]interface{})[j].(map[string]interface{})["id"] = socialConfig[i].ID
							body["socials"].([]interface{})[j].(map[string]interface{})["username"] = strings.Replace(body["socials"].([]interface{})[j].(map[string]interface{})["url"].(string), slicePath, "", -1)
							body["socials"].([]interface{})[j].(map[string]interface{})["url"] = slicePath + body["socials"].([]interface{})[j].(map[string]interface{})["username"].(string)

							body["socials"] = body["socials"].([]interface{})

						}
					}
				}
			}
		}

		_, err = users.InsertOne(context.Background(), bson.M{
			"id": middlewares.GenerateID(),
			"discord": bson.M{
				"id":            account,
				"username":      detailedAccount.Username,
				"discriminator": detailedAccount.Discriminator,
			},
			"avatar":     body["avatar"],
			"about":      body["about"],
			"url":        strings.ToLower(middlewares.RemoveSpecialChars(body["url"].(string))),
			"location":   body["location"],
			"occupation": body["occupation"],
			"birthday":   body["birthday"],
			"gender":     body["gender"],
			"language":   body["language"],
			"email":      body["email"],
			"likes":      []string{},
			"isVerified": false,
			"isPremium":  false,
			"socials":    body["socials"],
			"skills":     body["skills"],
			"roles":      body["roles"],
			"banner":     body["banner"],
			"pronouns":   body["pronouns"],
			"privacy": bson.M{
				"isShow":            body["privacy"] != nil && body["privacy"].(map[string]interface{})["isShow"] != nil && body["privacy"].(map[string]interface{})["isShow"] != "" && body["privacy"].(map[string]interface{})["isShow"] != "null",
				"isEmailPrivate":    body["privacy"] != nil && body["privacy"].(map[string]interface{})["isEmailPrivate"] != nil && body["privacy"].(map[string]interface{})["isEmailPrivate"] != "" && body["privacy"].(map[string]interface{})["isEmailPrivate"] != "null",
				"isGenderPrivate":   body["privacy"] != nil && body["privacy"].(map[string]interface{})["isGenderPrivate"] != nil && body["privacy"].(map[string]interface{})["isGenderPrivate"] != "" && body["privacy"].(map[string]interface{})["isGenderPrivate"] != "null",
				"isLocationPrivate": body["privacy"] != nil && body["privacy"].(map[string]interface{})["isLocationPrivate"] != nil && body["privacy"].(map[string]interface{})["isLocationPrivate"] != "" && body["privacy"].(map[string]interface{})["isLocationPrivate"] != "null",
				"isBirthdayPrivate": body["privacy"] != nil && body["privacy"].(map[string]interface{})["isBirthdayPrivate"] != nil && body["privacy"].(map[string]interface{})["isBirthdayPrivate"] != "" && body["privacy"].(map[string]interface{})["isBirthdayPrivate"] != "null",
			},
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
			"deletedAt": nil,
		})

		if err != nil {
			return c.Status(500).JSON(middlewares.NewError("Internal server error"))
		}
	}

	return c.JSON(middlewares.NewSuccess("User created", users))
}
