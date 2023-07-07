package biometryHandler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mmijangosFGE/validations-service/internal/core/ports"
	"github.com/mmijangosFGE/validations-service/pkg/functions"
	"github.com/mmijangosFGE/validations-service/pkg/messages"
	"github.com/mmijangosFGE/validations-service/pkg/requests"
	"github.com/mmijangosFGE/validations-service/pkg/responses"
)

// BiometryHandler - struct to create methods to establish communication
type BiometryHandler struct {
	biometryService ports.BiometryService
}

var _ ports.BiometryHandler = (*BiometryHandler)(nil)

// NewBiometryHandler - method to create a new instance of BiometryHandler
func NewHandler(biometryService ports.BiometryService) *BiometryHandler {
	return &BiometryHandler{
		biometryService: biometryService,
	}
}

// CompareFaces - method to compare faces
func (bh *BiometryHandler) CompareFaces(c *fiber.Ctx) error {
	// Get the request body params
	payload := new(requests.CompareFacesRequest)
	if err := c.BodyParser(payload); err != nil {
		response, ctx := responses.GeneralResponse(
			c,
			fiber.StatusBadRequest,
			false,
			messages.BadRequest,
		)
		return ctx.JSON(response)
	}
	// Payload to variables
	sourceImage := payload.SourceImage
	targetImage := payload.TargetImage
	similarityThreshold := payload.SimilarityThreshold
	// Validate if the pictures are not empty
	// and if the similarity threshold is between 0 and 1
	if strings.TrimSpace(sourceImage) == "" ||
		strings.TrimSpace(targetImage) == "" ||
		!functions.ValidateSimilarityThreshold(similarityThreshold) {
		response, ctx := responses.GeneralResponse(
			c,
			fiber.StatusBadRequest,
			false,
			messages.BadRequest,
		)
		return ctx.JSON(response)
	}
	// Call the service
	matched, status, err := bh.biometryService.CompareFaces(
		sourceImage,
		targetImage,
		similarityThreshold,
	)
	if err != nil {
		message := functions.GetServiceErrorMessage(status, err.Error())
		response, ctx := responses.GeneralResponse(
			c,
			status,
			false,
			message,
		)
		return ctx.JSON(response)
	}
	// Response is ok
	response, ctx := responses.GeneralResponse(
		c,
		status,
		matched,
		messages.Ok,
	)
	return ctx.JSON(response)
}
