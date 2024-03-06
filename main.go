package main

import "github.com/gofiber/fiber"

func aboutMessage(c *fiber.Ctx) {
	c.Send("You've reached Spotigram backend server, congrats!")
}

func setupRoutes(app *fiber.App) {
	app.Get("/about", aboutMessage)
}

func main() {
	app := fiber.New()
	setupRoutes(app)
	app.Listen(3000)
}
