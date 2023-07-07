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
	sourceImage string,
	targetImage string,
	similarityThreshold float64,
) (bool, int, error) {
	// Create a new session in the us-east-1 region.
	session, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"",
			),
		},
	})
	if err != nil {
		return false, fiber.StatusInternalServerError, err
	}
	// Create a service client.
	svc := rekognition.New(session)
	// validate if the pictures are valid urls
	if !functions.IsValidURL(sourceImage) || !functions.IsValidURL(targetImage) {
		return false, fiber.StatusInternalServerError, errors.New(messages.InvalidURL)
	}
	// Parse the images to byte arrays
	sourceImageBytes, err := functions.GetImageBytesFromURL(sourceImage)
	if err != nil {
		return false, fiber.StatusInternalServerError, err
	}
	targetImageBytes, err := functions.GetImageBytesFromURL(targetImage)
	if err != nil {
		return false, fiber.StatusInternalServerError, err
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
	result, err := svc.CompareFaces(input)
	if err != nil {
		return false, fiber.StatusInternalServerError, err
	}
	// Validate if the faces are similar
	if len(result.FaceMatches) > 0 {
		return true, fiber.StatusOK, nil
	}
	// If the faces are not similar, return false but not an error
	return false, fiber.StatusOK, nil
}
