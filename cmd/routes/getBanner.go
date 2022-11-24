package routes

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.dbio.me/cmd/configuration"
	"go.dbio.me/cmd/middlewares"
	"go.dbio.me/cmd/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetBanner(c *fiber.Ctx) error {
	DConfig := configuration.GetConfig()
	var collectionData = DConfig.Collection

	var db *mongo.Client = c.Locals("db").(*mongo.Client)
	users := db.Database("dbio").Collection(collectionData)

	var user types.Entity

	err := users.FindOne(context.Background(), bson.M{"url": c.Params("id")}).Decode(&user)

	if err != nil {
		return c.Status(404).JSON(middlewares.NewError("User not found"))
	}

	if user.Banner == "" {
		return c.Status(404).JSON(middlewares.NewError("Banner not found"))
	}

	url := user.Banner

	resp, err := http.Get(url)
	if err != nil {
		return c.Status(404).JSON(middlewares.NewError("Banner not found"))
	}

	image, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(404).JSON(middlewares.NewError("Banner not found"))
	}

	c.Set("Content-Type", "image/png")

	_, err = c.Write(image)
	if err != nil {
		return c.Status(404).JSON(middlewares.NewError("Banner not found"))
	}

	return nil
}
