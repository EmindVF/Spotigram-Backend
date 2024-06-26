package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// A handler to send about message.
func AboutHandler(ctx *fiber.Ctx) error {
	ctx.SendString("You have reached a test version of Spotigram!")
	return nil
}

// A handler to send not found message.
func NotFoundHandler(ctx *fiber.Ctx) error {
	ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status": "fail",
		"message": fmt.Sprintf(
			"path \"%v\" with method \"%v\" does not exist on this server",
			ctx.Path(),
			ctx.Method()),
	})
	return nil
}
