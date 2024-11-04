package jigsawstack

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	vOCREndpoint            Endpoint = "/v1/vocr"
	vObjectEndpoint         Endpoint = "/v1/ai/object_detection"
	imageGenerationEndpoint Endpoint = "/v1/ai/image_generation"
)

type (
	// model
	// string
	// default: "sdxl"
	//
	// The model to use for the generation. Default is sdxl
	//
	//	sd1.5 - Stable Diffusion v1.5
	//	sdxl - Stable Diffusion XL
	//	ead1.0 - Anime Diffusion
	//	rv1.3 - Realistic Vision v1.3
	//	rv3 - Realistic Vision v3
	//	rv5.1 - Realistic Vision v5.1
	//	ar1.8 - AbsoluteReality v1.8.1

	// ImageGenerationRequest represents a request structure for image
	// generation API.
	ImageGenerationRequest struct {
		Prompt string `json:"prompt"`
		Model  string `json:"model,omitempty"`
		Size   string `json:"size"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}
	// ImageGenerationResponse represents a response structure for image
	// generation API.
	ImageGenerationResponse struct {
		Success bool   `json:"success"`
		Image   string `json:"image"`
	}
	// visionRequest represents a request structure for VOCR API.
	visionRequest struct {
		// Prompt is the prompt used in ocr. If the request is for
		// object detection, this field is not required.
		Prompt string `json:"prompt,omitempty"`
		// URL is the url of the image to use as the image.
		//
		// Not required if the StoreKey is not provided.
		URL string `json:"image_url"`
		// Key is the key of the file to use as the image.
		//
		// Not required if the ImageURL is not provided.
		Key string `json:"file_store_key"`
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

// VCOROption is the option for VOCR.
type VCOROption func(*visionRequest)

// WithKey sets the key of the file to use as the image.
func WithKey(key string) VCOROption {
	return func(params *visionRequest) { params.Key = key }
}

// WithURL sets the URL of the image to use as the image.
func WithURL(url string) VCOROption {
	return func(params *visionRequest) { params.URL = url }
}

// VOCR performs a visual object recognition (VOCR) task on an image.
//
// POST https://api.jigsawstack.com/v1/vocr
//
// https://docs.jigsawstac.com/api-reference/ai/vision
func (j *JigsawStack) VOCR(
	ctx context.Context,
	prompt string,
	opt VCOROption,
) (string, error) {
	params := visionRequest{
		Prompt: prompt,
	}
	opt(&params)
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
	params visionRequest,
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

// ImageGeneration generates an image from a prompt and parameters.
func (j *JigsawStack) ImageGeneration(
	ctx context.Context,
	request ImageGenerationRequest,
) (response ImageGenerationResponse, err error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+string(imageGenerationEndpoint),
		builders.WithBody(request),
	)
	if err != nil {
		return
	}
	var resp ImageGenerationResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}
