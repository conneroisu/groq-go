package jigsawstack

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	imageGenerationEndpoint = "v1/ai/image_generation"
)

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
type (
	ImageGenerationRequest struct {
		Prompt string `json:"prompt"`
		Model  string `json:"model,omitempty"`
		Size   string `json:"size"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}
	ImageGenerationResponse struct {
		Success bool   `json:"success"`
		Image   string `json:"image"`
	}
)

// Default is sdxl
func (j *JigsawStack) ImageGeneration(
	ctx context.Context,
	request ImageGenerationRequest,
) (response ImageGenerationResponse, err error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+imageGenerationEndpoint,
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
