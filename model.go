// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package groq

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/conneroisu/groq-go/internal/apijson"
	"github.com/conneroisu/groq-go/internal/requestconfig"
	"github.com/conneroisu/groq-go/option"
)

// ModelService contains methods and other services that help with interacting with
// the groq API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewModelService] method instead.
type ModelService struct {
	Options []option.RequestOption
}

// NewModelService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewModelService(opts ...option.RequestOption) (r *ModelService) {
	r = &ModelService{}
	r.Options = opts
	return
}

// Get a specific model
func (r *ModelService) Get(ctx context.Context, model string, opts ...option.RequestOption) (res *Model, err error) {
	opts = append(r.Options[:], opts...)
	if model == "" {
		err = errors.New("missing required model parameter")
		return
	}
	path := fmt.Sprintf("openai/v1/models/%s", model)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// get all available models
func (r *ModelService) List(ctx context.Context, opts ...option.RequestOption) (res *ModelListResponse, err error) {
	opts = append(r.Options[:], opts...)
	path := "openai/v1/models"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Delete a model
func (r *ModelService) Delete(ctx context.Context, model string, opts ...option.RequestOption) (res *ModelDeleteResponse, err error) {
	opts = append(r.Options[:], opts...)
	if model == "" {
		err = errors.New("missing required model parameter")
		return
	}
	path := fmt.Sprintf("openai/v1/models/%s", model)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Describes an OpenAI model offering that can be used with the API.
type Model struct {
	// The model identifier, which can be referenced in the API endpoints.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) when the model was created.
	Created int64 `json:"created,required"`
	// The object type, which is always "model".
	Object ModelObject `json:"object,required"`
	// The organization that owns the model.
	OwnedBy string    `json:"owned_by,required"`
	JSON    modelJSON `json:"-"`
}

// modelJSON contains the JSON metadata for the struct [Model]
type modelJSON struct {
	ID          apijson.Field
	Created     apijson.Field
	Object      apijson.Field
	OwnedBy     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Model) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelJSON) RawJSON() string {
	return r.raw
}

// The object type, which is always "model".
type ModelObject string

const (
	ModelObjectModel ModelObject = "model"
)

func (r ModelObject) IsKnown() bool {
	switch r {
	case ModelObjectModel:
		return true
	}
	return false
}

type ModelListResponse struct {
	Data   []Model                 `json:"data,required"`
	Object ModelListResponseObject `json:"object,required"`
	JSON   modelListResponseJSON   `json:"-"`
}

// modelListResponseJSON contains the JSON metadata for the struct
// [ModelListResponse]
type modelListResponseJSON struct {
	Data        apijson.Field
	Object      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ModelListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelListResponseJSON) RawJSON() string {
	return r.raw
}

type ModelListResponseObject string

const (
	ModelListResponseObjectList ModelListResponseObject = "list"
)

func (r ModelListResponseObject) IsKnown() bool {
	switch r {
	case ModelListResponseObjectList:
		return true
	}
	return false
}

type ModelDeleteResponse struct {
	ID      string                  `json:"id,required"`
	Deleted bool                    `json:"deleted,required"`
	Object  string                  `json:"object,required"`
	JSON    modelDeleteResponseJSON `json:"-"`
}

// modelDeleteResponseJSON contains the JSON metadata for the struct
// [ModelDeleteResponse]
type modelDeleteResponseJSON struct {
	ID          apijson.Field
	Deleted     apijson.Field
	Object      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ModelDeleteResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelDeleteResponseJSON) RawJSON() string {
	return r.raw
}
