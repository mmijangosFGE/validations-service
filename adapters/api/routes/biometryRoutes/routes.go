package biometryRoutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mmijangosFGE/validations-service/internal/core/services/biometryService"
	"github.com/mmijangosFGE/validations-service/internal/handlers/biometryHandler"
)

// InitRoutes - function to initialize routes
func InitRoutes(
	app fiber.Router,
	basicAuth fiber.Handler,
) {
	// create a prefix to routes
	v1 := app.Group("/v1")
	// Instance of biometry service
	biometryService := biometryService.NewService()
	// Instance of biometry handler
	biometryHandler := biometryHandler.NewHandler(biometryService)
	// Create a route to compare faces
	biometryRoutes := v1.Group("/biometry")
	// Create the routes
	biometryRoutes.Post("/compare-faces", biometryHandler.CompareFaces)

}
