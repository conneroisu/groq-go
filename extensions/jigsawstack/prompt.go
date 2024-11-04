package jigsawstack

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	promptCreateEndpoint Endpoint = "/v1/prompt_engine"
)

type (
	// PromptRunResponse represents a response structure for prompt run API.
	PromptRunResponse struct {
		Success bool   `json:"success"`
		Result  string `json:"result"`
	}
	// PromptCreateInput represents an entry in a prompt create request.
	PromptCreateInput struct {
		Key          string `json:"key"`
		Optional     bool   `json:"optional"`
		InitialValue string `json:"initial_value"`
	}
	// PromptCreateRequest represents a request structure for prompt create API.
	PromptCreateRequest struct {
		Prompt       string              `json:"prompt"`
		Inputs       []PromptCreateInput `json:"inputs"`
		ReturnPrompt string              `json:"return_prompt"`
		PromptGuard  []string            `json:"prompt_guard"`
		Optimize     bool                `json:"optimize_prompt,omitempty"`
		UseInternet  bool                `json:"use_internet,omitempty"`
	}
	// PromptResponse represents a response structure for prompt create API.
	PromptResponse struct {
		Success bool   `json:"success"`
		ID      string `json:"prompt_engine_id"`
	}
	// PromptEngine represents a prompt engine.
	PromptEngine struct {
		ID           string    `json:"id"`
		Prompt       string    `json:"prompt"`
		Inputs       any       `json:"inputs"`
		ReturnPrompt any       `json:"return_prompt"`
		CreatedAt    time.Time `json:"created_at"`
	}
	// PromptListResponse represents a response structure for prompt list API.
	PromptListResponse struct {
		Success       bool           `json:"success"`
		PromptEngines []PromptEngine `json:"prompt_engines"`
		Page          int            `json:"page"`
		Limit         int            `json:"limit"`
		HasMore       bool           `json:"has_more"`
	}
)

// PromptGet gets a specific prompt.
//
// GET https://api.jigsawstack.com/v1/prompt_engine/{id}
//
// https://docs.jigsawstack.com/api-reference/prompt-engine/retrieve
func (j *JigsawStack) PromptGet(
	ctx context.Context,
	id string,
) (response PromptEngine, err error) {
	uri := j.baseURL + string(promptCreateEndpoint) + "/" + id
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodGet,
		uri,
	)
	if err != nil {
		return
	}
	var resp PromptEngine
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}

// PromptList lists prompts.
//
// https://docs.jigsawstack.com/api-reference/prompt-engine/list
//
// GET https://api.jigsawstack.com/v1/prompt_engine
func (j *JigsawStack) PromptList(
	ctx context.Context,
	page int,
	limit int,
) (response PromptListResponse, err error) {
	uri := j.baseURL + string(promptCreateEndpoint) + "?page=" + strconv.Itoa(page) + "&limit=" + strconv.Itoa(limit)
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodGet,
		uri,
	)
	if err != nil {
		return
	}
	var resp PromptListResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}

// PromptCreate creates a prompt.
//
// https://docs.jigsawstack.com/api-reference/prompt-engine/create
//
// POST https://api.jigsawstack.com/v1/prompt_engine
func (j *JigsawStack) PromptCreate(
	ctx context.Context,
	request PromptCreateRequest,
) (response PromptResponse, err error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+string(promptCreateEndpoint),
		builders.WithBody(request),
	)
	if err != nil {
		return
	}
	var resp PromptResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}

// PromptDelete deletes a specific prompt.
//
// https://api.jigsawstack.com/v1/prompt_engine/{id}
//
// https://docs.jigsawstack.com/api-reference/prompt-engine/delete
func (j *JigsawStack) PromptDelete(
	ctx context.Context,
	id string,
) (response PromptResponse, err error) {
	// TODO: may need to sanitize the id
	uri := j.baseURL + string(promptCreateEndpoint) + "/" + id
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodDelete,
		uri,
	)
	if err != nil {
		return
	}
	var resp PromptResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}

// PromptRun runs a specific prompt with the given inputs.
//
// https://api.jigsawstack.com/v1/prompt_engine/{id}
//
// https://docs.jigsawstack.com/api-reference/prompt-engine/run
func (j *JigsawStack) PromptRun(
	ctx context.Context,
	id string,
	inputs map[string]any,
) (response PromptRunResponse, err error) {
	// TODO: may need to sanitize the id
	uri := j.baseURL + string(promptCreateEndpoint) + "/" + id
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		uri,
		builders.WithBody(inputs),
	)
	if err != nil {
		return
	}
	var resp PromptRunResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}

// PromptRunDirect runs new prompt with the given inputs.
//
// https://docs.jigsawstack.com/api-reference/prompt-engine/run-direct
//
// https://api.jigsawstack.com/v1/prompt_engine/run
func (j *JigsawStack) PromptRunDirect(
	ctx context.Context,
	request PromptCreateRequest,
	inputs map[string]any,
) (response PromptRunResponse, err error) {
	type combinedRequest struct {
		PromptCreateRequest
		Inputs map[string]any `json:"inputs"`
	}
	var combinedReq combinedRequest
	combinedReq.PromptCreateRequest = request
	combinedReq.Inputs = inputs
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+string(promptCreateEndpoint),
		builders.WithBody(combinedReq),
	)
	if err != nil {
		return
	}
	var resp PromptRunResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}
