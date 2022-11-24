package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.dbio.me/cmd/middlewares"
)

func GetLanguages(c *fiber.Ctx) error {

	var allLanguages []string

	allLanguages =
		[]string{
			"English",
			"Spanish",
			"French",
			"German",
			"Russian",
			"Japanese",
			"Chinese",
			"Korean",
			"Arabic",
			"Hindi",
			"Portuguese",
			"Indonesian",
			"Turkish",
			"Persian",
			"Malay",
			"Thai",
			"Vietnamese",
			"Urdu",
			"Romanian",
			"Bengali",
			"Polish",
			"Ukrainian",
			"Tagalog",
			"Serbian",
			"Swedish",
			"Kannada",
			"Malayalam",
			"Kurdish",
			"Marathi",
			"Kazakh",
			"Kinyarwanda",
		}

	return c.Status(200).JSON(middlewares.NewSuccess("Languages fetched", allLanguages))
}
