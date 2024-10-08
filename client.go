package groq

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

//go:generate go run ./scripts/generate-models/
//go:generate make docs

// Format is the format of a response.
// string
type (
	Format string
	// Client is a Groq api client.
	Client struct {
		groqAPIKey         string // Groq API key
		orgID              string // OrgID is the organization ID for the client.
		baseURL            string // Base URL for the client.
		emptyMessagesLimit uint   // EmptyMessagesLimit is the limit for the empty messages.
		requestBuilder     requestBuilder
		requestFormBuilder formBuilder
		createFormBuilder  func(body io.Writer) formBuilder

		client *http.Client // Client is the HTTP client to use
		logger *slog.Logger // Logger is the logger for the client.
	}

	// RateLimitHeaders struct represents Groq rate limits headers.
	RateLimitHeaders struct {
		LimitRequests     int       `json:"x-ratelimit-limit-requests"`     // LimitRequests is the limit requests of the rate limit headers.
		LimitTokens       int       `json:"x-ratelimit-limit-tokens"`       // LimitTokens is the limit tokens of the rate limit headers.
		RemainingRequests int       `json:"x-ratelimit-remaining-requests"` // RemainingRequests is the remaining requests of the rate limit headers.
		RemainingTokens   int       `json:"x-ratelimit-remaining-tokens"`   // RemainingTokens is the remaining tokens of the rate limit headers.
		ResetRequests     ResetTime `json:"x-ratelimit-reset-requests"`     // ResetRequests is the reset requests of the rate limit headers.
		ResetTokens       ResetTime `json:"x-ratelimit-reset-tokens"`       // ResetTokens is the reset tokens of the rate limit headers.
	}
	// ResetTime is a time.Time wrapper for the rate limit reset time.
	// string
	ResetTime string

	// Opts is a function that sets options for a Groq client.
	Opts func(*Client)

	requestOption  func(*requestOptions)
	fullURLOptions struct {
		model string
	}
	fullURLOption func(*fullURLOptions)

	// Usage Represents the total token usage per request to Groq.
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}

	requestOptions struct {
		body   any
		header http.Header
	}
)

const (
	// FormatText is the text format. It is the default format of a
	// response.
	FormatText Format = "text"
	// FormatJSON is the JSON format. There is no support for streaming with
	// JSON format selected.
	FormatJSON Format = "json"

	// groqAPIURLv1 is the base URL for the Groq API.
	groqAPIURLv1 = "https://api.groq.com/openai/v1"
)

// NewClient creates a new Groq client.
func NewClient(groqAPIKey string, opts ...Opts) (*Client, error) {
	if groqAPIKey == "" {
		return nil, fmt.Errorf("groq api key is required")
	}
	c := &Client{
		groqAPIKey:         groqAPIKey,
		client:             http.DefaultClient,
		logger:             slog.Default(),
		baseURL:            groqAPIURLv1,
		emptyMessagesLimit: 10,
		createFormBuilder: func(body io.Writer) formBuilder {
			return newFormBuilder(body)
		},
		requestBuilder: newRequestBuilder(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// fullURL returns full URL for request.
func (c *Client) fullURL(suffix Endpoint, setters ...fullURLOption) string {
	baseURL := strings.TrimRight(c.baseURL, "/")
	args := fullURLOptions{}
	for _, setter := range setters {
		setter(&args)
	}
	return fmt.Sprintf("%s%s", baseURL, suffix)
}

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
func WithLogger(logger *slog.Logger) Opts {
	return func(c *Client) {
		c.logger = logger
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
		header: http.Header{},
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

// response is an interface for a response.
type response interface {
	SetHeader(http.Header)
}

func (c *Client) sendRequest(req *http.Request, v response) error {
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

func sendRequestStream[T streamer](
	client *Client,
	req *http.Request,
) (*streamReader[T], error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.client.Do(
		req,
	) //nolint:bodyclose // body is closed in stream.Close()
	if err != nil {
		return new(streamReader[T]), err
	}
	if isFailureStatusCode(resp) {
		return new(streamReader[T]), client.handleErrorResp(resp)
	}
	return &streamReader[T]{
		emptyMessagesLimit: client.emptyMessagesLimit,
		reader:             bufio.NewReader(resp.Body),
		response:           resp,
		errAccumulator:     newErrorAccumulator(),
		Header:             resp.Header,
	}, nil
}

func (c *Client) setCommonHeaders(req *http.Request) {
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

func withModel(model model) fullURLOption {
	return func(args *fullURLOptions) {
		args.model = string(model)
	}
}

func (c *Client) handleErrorResp(resp *http.Response) error {
	var errRes errorResponse
	err := json.NewDecoder(resp.Body).Decode(&errRes)
	if err != nil || errRes.Error == nil {
		reqErr := &requestError{
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

// String returns the string representation of the ResetTime.
func (r ResetTime) String() string {
	return string(r)
}

// Time returns the time.Time representation of the ResetTime.
func (r ResetTime) Time() time.Time {
	d, _ := time.ParseDuration(string(r))
	return time.Now().Add(d)
}
