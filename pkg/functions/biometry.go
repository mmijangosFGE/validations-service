package functions

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mmijangosFGE/validations-service/pkg/messages"
)

func IsValidURL(input string) bool {
	_, err := url.ParseRequestURI(input)
	return err == nil
}

func GetImageBytesFromURL(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	imageBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}

func ValidateSimilarityThreshold(similarityThreshold float64) bool {
	return similarityThreshold >= 0 && similarityThreshold <= 1
}

func GetServiceErrorMessage(
	httpStatus int,
	err string,
) string {
	// Validate if error exists and error type
	if httpStatus == fiber.StatusOK {
		return err
	}
	// Set error log
	errorLogger(err)
	// Return error message
	return messages.InternalServerError
}

func errorLogger(err string) {
	// config logger to set date and time
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stderr)
	// set error message
	errLog := errors.New(err)
	// visualize error message
	log.Printf("Error: %v", errLog)
}
