package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.dbio.me/cmd/configuration"
	"go.dbio.me/cmd/middlewares"
)

func GetSocials(c *fiber.Ctx) error {
	config := configuration.GetSocials()

	return c.Status(200).JSON(middlewares.NewSuccess("Successfully fetched socials", config))
}

func GetSortings(c *fiber.Ctx) error {
	config := configuration.GetSort()

	return c.Status(200).JSON(middlewares.NewSuccess("Successfully fetched sortings", config))
}
