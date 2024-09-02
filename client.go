package groq

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

const (
	groqAPIURLv1 = "https://api.groq.com/openai/v1"
)

// Client is a Groq api client.
type Client struct {
	groqAPIKey         string         // Groq API key
	orgID              string         // OrgID is the organization ID for the client.
	baseURL            string         // Base URL for the client.
	client             *http.Client   // Client is the HTTP client to use
	models             *ModelResponse // Models is the list of models available to the client.
	requestBuilder     RequestBuilder
	requestFormBuilder FormBuilder
	logger             zerolog.Logger // Logger is the logger for the client.
}

// Contains returns true if the model is in the list of models.
func (m *ModelResponse) contains(model string) bool {
	for _, m := range m.Data {
		if m.ID == model {
			return true
		}
	}
	return false
}

// NewClient creates a new Groq client.
func NewClient(groqAPIKey string, opts ...Opts) (*Client, error) {
	c := &Client{
		groqAPIKey: groqAPIKey,
		client:     http.DefaultClient,
		logger: zerolog.New(os.Stderr).
			Level(zerolog.DebugLevel).
			With().
			Timestamp().
			Logger(),
		baseURL: groqAPIURLv1,
	}
	err := c.GetModels()
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// GetModels gets the list of models from the Groq API.
func (c *Client) GetModels() error {
	req, err := http.NewRequest("GET", c.baseURL+"/models", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.groqAPIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var modelsResponse ModelResponse
	err = json.Unmarshal(bodyText, &modelsResponse)
	if err != nil {
		return err
	}
	c.models = &modelsResponse
	return nil
}

// ModelResponse is a response from the models endpoint.
type ModelResponse struct {
	Object string `json:"object"`
	Data   []struct {
		ID             string `json:"id"`
		Object         string `json:"object"`
		Created        int    `json:"created"`
		OwnedBy        string `json:"owned_by"`
		Active         bool   `json:"active"`
		ContextWindow  int    `json:"context_window,omitempty"`
		PublicApps     any    `json:"public_apps"`
		ContextWindow0 int    `json:"context_ window,omitempty"`
	} `json:"data"`
}

// Opts is a function that sets options for a Groq client.
type Opts func(*Client)

// WithClient sets the client for the Groq client.
func WithClient(client *http.Client) Opts {
	return func(c *Client) {
		c.client = client
	}
}

// WithBaseURL sets the base URL for the Groq client.
func WithBaseURL(baseURL string) Opts {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithLogger sets the logger for the Groq client.
func WithLogger(logger zerolog.Logger) Opts {
	return func(c *Client) {
		c.logger = logger
	}
}

// Usage Represents the total token usage per request to Groq.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
type requestOptions struct {
	body   any
	header http.Header
}

type requestOption func(*requestOptions)

func withBody(body any) requestOption {
	return func(args *requestOptions) {
		args.body = body
	}
}

func withContentType(contentType string) requestOption {
	return func(args *requestOptions) {
		args.header.Set("Content-Type", contentType)
	}
}

func (c *Client) newRequest(
	ctx context.Context,
	method, url string,
	setters ...requestOption,
) (*http.Request, error) {
	// Default Options
	args := &requestOptions{
		body:   nil,
		header: make(http.Header),
	}
	for _, setter := range setters {
		setter(args)
	}
	req, err := c.requestBuilder.Build(
		ctx,
		method,
		url,
		args.body,
		args.header,
	)
	if err != nil {
		return nil, err
	}
	c.setCommonHeaders(req)
	return req, nil
}

// Response is an interface for a response.
type Response interface {
	SetHeader(http.Header)
}

func (c *Client) sendRequest(req *http.Request, v Response) error {
	req.Header.Set("Accept", "application/json")

	// Check whether Content-Type is already set, Upload Files API requires
	// Content-Type == multipart/form-data
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if v != nil {
		v.SetHeader(res.Header)
	}

	if isFailureStatusCode(res) {
		return c.handleErrorResp(res)
	}

	return decodeResponse(res.Body, v)
}

func (c *Client) sendRequestRaw(
	req *http.Request,
) (response RawResponse, err error) {
	resp, err := c.config.HTTPClient.Do(
		req,
	) //nolint:bodyclose // body should be closed by outer function
	if err != nil {
		return
	}

	if isFailureStatusCode(resp) {
		err = c.handleErrorResp(resp)
		return
	}

	response.SetHeader(resp.Header)
	response.ReadCloser = resp.Body
	return
}

func sendRequestStream[T streamable](
	client *Client,
	req *http.Request,
) (*streamReader[T], error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.config.HTTPClient.Do(
		req,
	) //nolint:bodyclose // body is closed in stream.Close()
	if err != nil {
		return new(streamReader[T]), err
	}
	if isFailureStatusCode(resp) {
		return new(streamReader[T]), client.handleErrorResp(resp)
	}
	return &streamReader[T]{
		emptyMessagesLimit: client.EmptyMessagesLimit,
		reader:             bufio.NewReader(resp.Body),
		response:           resp,
		errAccumulator:     NewErrorAccumulator(),
		unmarshaler:        &JSONUnmarshaler{},
		Header:             resp.Header,
	}, nil
}

func (c *Client) setCommonHeaders(req *http.Request) {
	// https://learn.microsoft.com/en-us/azure/cognitive-services/openai/reference#authentication
	// Azure API Key authentication
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.groqAPIKey))
	if c.orgID != "" {
		req.Header.Set("OpenAI-Organization", c.orgID)
	}
}

func isFailureStatusCode(resp *http.Response) bool {
	return resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusBadRequest
}

func decodeResponse(body io.Reader, v any) error {
	if v == nil {
		return nil
	}

	switch o := v.(type) {
	case *string:
		return decodeString(body, o)
	case *audioTextResponse:
		return decodeString(body, &o.Text)
	default:
		return json.NewDecoder(body).Decode(v)
	}
}

func decodeString(body io.Reader, output *string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	*output = string(b)
	return nil
}

type fullURLOptions struct {
	model string
}

type fullURLOption func(*fullURLOptions)

func withModel(model string) fullURLOption {
	return func(args *fullURLOptions) {
		args.model = model
	}
}

func (c *Client) handleErrorResp(resp *http.Response) error {
	var errRes ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errRes)
	if err != nil || errRes.Error == nil {
		reqErr := &RequestError{
			HTTPStatusCode: resp.StatusCode,
			Err:            err,
		}
		if errRes.Error != nil {
			reqErr.Err = errRes.Error
		}
		return reqErr
	}

	errRes.Error.HTTPStatusCode = resp.StatusCode
	return errRes.Error
}
