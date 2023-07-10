package ports

import (
	"github.com/gofiber/fiber/v2"
)

// BiometryHandler - interface to create methods to establish communication
// between the requests and the service
type BiometryHandler interface {
	CompareFaces(c *fiber.Ctx) error
}

// BiometryService - interface to create methods to establish communication
// between the handlers and the service
type BiometryService interface {
	CompareFaces(
		SimilarityThreshold float64,
		SourceImage string,
		TargetImage string,
	) (bool, int, error)
}
