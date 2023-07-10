package biometryService

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/gofiber/fiber/v2"
	"github.com/mmijangosFGE/validations-service/pkg/functions"
	"github.com/mmijangosFGE/validations-service/pkg/messages"
)

// BiometryService - struct of biometry service
type BiometryService struct{}

func NewService() *BiometryService {
	return &BiometryService{}
}

// CompareFaces - method to compare faces
func (s *BiometryService) CompareFaces(
	similarityThreshold float64,
	sourceImage string,
	targetImage string,
) (bool, int, error) {
	// Create a new session in the us-east-1 region.
	session, errSession := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"",
			),
		},
	})
	if errSession != nil {
		return false, fiber.StatusInternalServerError, errSession
	}
	// Create a service client.
	svc := rekognition.New(session)
	// validate if the pictures are valid urls
	if !functions.IsValidURL(sourceImage) || !functions.IsValidURL(targetImage) {
		return false,
			fiber.StatusInternalServerError, errors.New(messages.InvalidURL)
	}
	// Parse the images to byte arrays
	sourceImageBytes, errSourceImage := functions.GetImageBytesFromURL(sourceImage)
	if errSourceImage != nil {
		return false, fiber.StatusInternalServerError, errSession
	}
	targetImageBytes, errTargetImage := functions.GetImageBytesFromURL(targetImage)
	if errTargetImage != nil {
		return false, fiber.StatusInternalServerError, errSession
	}
	// Create the input for the request
	input := &rekognition.CompareFacesInput{
		SimilarityThreshold: aws.Float64(similarityThreshold),
		SourceImage: &rekognition.Image{
			Bytes: sourceImageBytes,
		},
		TargetImage: &rekognition.Image{
			Bytes: targetImageBytes,
		},
	}
	// Call operation
	result, errCompare := svc.CompareFaces(input)
	if errCompare != nil {
		return false, fiber.StatusInternalServerError, errSession
	}
	// Validate if the faces are similar
	if len(result.FaceMatches) > 0 {
		return true, fiber.StatusOK, nil
	}
	// If the faces are not similar, return false but not an error
	return false, fiber.StatusOK, nil
}
