package responses

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mmijangosFGE/validations-service/pkg/constants"
)

/*
GeneralResponse -
Create json to response to client without data
*/
func GeneralResponse(
	c *fiber.Ctx,
	httpStatus int,
	success bool,
	message string,
) (fiber.Map, *fiber.Ctx) {
	// Set header to response with json
	c.Set("Content-Type", "application/json")
	// Set http status to response
	if err := c.SendStatus(httpStatus); err != nil {
		panic(err)
	}
	// Build a json response
	json := fiber.Map{
		"success": success,
		"message": message,
	}
	// return json and context
	return json, c
}

/*
ObjectResponse -
Create json to response to client with object data
*/
func ObjectResponse(
	c *fiber.Ctx,
	httpStatus int,
	success bool,
	message string,
	data fiber.Map,
) (fiber.Map, *fiber.Ctx) {
	// Set header to response with json
	c.Set(constants.HeaderContentType, constants.AppJsonHeader)
	// Set http status to response
	if err := c.SendStatus(httpStatus); err != nil {
		panic(err)
	}
	// Build a json response
	json := fiber.Map{
		"success": success,
		"message": message,
		"data":    data,
	}
	// return json and context
	return json, c
}
