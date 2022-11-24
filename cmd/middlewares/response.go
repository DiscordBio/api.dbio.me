package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func NewError(msg string, data ...interface{}) fiber.Map {
	return fiber.Map{
		"success": false,
		"error":   msg,
		"data":    data,
	}
}

func NewSuccess(msg string, data interface{}) fiber.Map {
	return fiber.Map{
		"success": true,
		"message": msg,
		"data":    data,
	}
}
