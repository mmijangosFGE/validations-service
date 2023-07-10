package requests

// CompareFacesRequest - struct of request to compare faces
type CompareFacesRequest struct {
	SimilarityThreshold float64 `json:"similarityThreshold" validate:"" default:"0.9"`
	SourceImage         string  `json:"sourceImage" validate:"required"`
	TargetImage         string  `json:"targetImage" validate:"required"`
}
