package jigsawstack

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	vOCREndpoint    Endpoint = "/v1/vocr"
	vObjectEndpoint Endpoint = "/v1/ai/object_detection"
)

type (
	// VisionRequest represents a request structure for VOCR API.
	VisionRequest struct {
		// Prompt is the prompt used in ocr. If the request is for
		// object detection, this field is not required.
		Prompt string `json:"prompt,omitempty"`
		// ImageURL is the url of the image to use as the image.
		//
		// Not required if the StoreKey is not provided.
		ImageURL string `json:"image_url"`
		// StoreKey is the key of the file to use as the image.
		//
		// Not required if the ImageURL is not provided.
		StoreKey string `json:"file_store_key"`
	}
	// VOCRResponse represents a response structure for VOCR API.
	VOCRResponse struct {
		// Success is a boolean indicating whether the request was
		// successful.
		Success bool `json:"success"`
		// Context is the context of the image.
		Context string `json:"context"`
		// Width is the width of the image.
		Width int `json:"width"`
		// Height is the height of the image.
		Height int `json:"height"`
		// Tags is a list of tags detected in the image.
		Tags []string `json:"tags"`
		// HasText is a boolean indicating whether the image contains
		// text.
		HasText bool `json:"has_text"`
		// Sections is a list of sections detected in the image.
		Sections []any `json:"sections"`
	}
	// VisionObjectResponse represents a response structure for VOD API.
	VisionObjectResponse struct {
		// Success is a boolean indicating whether the request was
		Success bool `json:"success"`
		// Width is the width of the image.
		Width int `json:"width"`
		// Height is the height of the image.
		Height int `json:"height"`
		// Tags is a list of tags detected in the image.
		Tags []string `json:"tags"`
		// Objects is a list of objects detected in the image.
		Objects []struct {
			Name       string  `json:"name"`
			Confidence float64 `json:"confidence"`
			Bounds     struct {
				TopLeft struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"top_left"`
				TopRight struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"top_right"`
				BottomRight struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"bottom_right"`
				BottomLeft struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"bottom_left"`
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"bounds"`
		} `json:"objects"`
	}
)

// VOCR performs a visual object recognition (VOCR) task on an image.
//
// POST https://api.jigsawstack.com/v1/vocr
//
// https://docs.jigsawstac.com/api-reference/ai/vision
func (j *JigsawStack) VOCR(
	ctx context.Context,
	params VisionRequest,
) (string, error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+string(vOCREndpoint),
		builders.WithBody(params),
	)
	if err != nil {
		return "", err
	}
	var resp VOCRResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return "", err
	}
	return "", nil
}

// VisionObjectDetection performs a visual object detection (VOD) task on an
// image.
//
// POST https://api.jigsawstack.com/v1/ai/object_detection
//
// https://docs.jigsawstack.com/api-reference/ai/object-detection
func (j *JigsawStack) VisionObjectDetection(
	ctx context.Context,
	params VisionRequest,
) (string, error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+string(vObjectEndpoint),
		builders.WithBody(params),
	)
	if err != nil {
		return "", err
	}
	var resp VisionObjectResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return "", err
	}
	return "", nil
}
