package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mmijangosFGE/validations-service/adapters/api/routes/biometryRoutes"
)

// BindRoutes - function to bind routes
func BindRoutes(
	app *fiber.App,
	basicAuth *fiber.Handler,
) {
	// Set the home route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	// Init the biometry routes
	biometryRoutes.InitRoutes(app, *basicAuth)
}
