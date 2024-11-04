package groq

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/conneroisu/groq-go/pkg/groqerr"
	"github.com/conneroisu/groq-go/pkg/models"
	"github.com/conneroisu/groq-go/pkg/streams"
)

//go:generate go run ./scripts/generate-models/
//go:generate go run github.com/princjef/gomarkdoc/cmd/gomarkdoc@v1.1.0 -o README.md -e .

type (
	// Client is a Groq api client.
	Client struct {
		// Groq API key
		groqAPIKey string
		// OrgID is the organization ID for the client.
		orgID string
		// Base URL for the client.
		baseURL string
		// EmptyMessagesLimit is the limit for the empty messages.
		emptyMessagesLimit uint

		header             builders.Header
		requestFormBuilder builders.FormBuilder

		// Client is the HTTP client to use
		client *http.Client
		// Logger is the logger for the client.
		logger *slog.Logger
	}
	// Opts is a function that sets options for a Groq client.
	Opts func(*Client)
)

// WithClient sets the client for the Groq client.
func WithClient(client *http.Client) Opts {
	return func(c *Client) { c.client = client }
}

// WithBaseURL sets the base URL for the Groq client.
func WithBaseURL(baseURL string) Opts {
	return func(c *Client) { c.baseURL = baseURL }
}

// WithLogger sets the logger for the Groq client.
func WithLogger(logger *slog.Logger) Opts {
	return func(c *Client) { c.logger = logger }
}

type (
	// Format is the format of a response.
	// string
	Format string
	// RateLimitHeaders struct represents Groq rate limits headers.
	RateLimitHeaders struct {
		// LimitRequests is the limit requests of the rate limit
		// headers.
		LimitRequests int `json:"x-ratelimit-limit-requests"`
		// LimitTokens is the limit tokens of the rate limit headers.
		LimitTokens int `json:"x-ratelimit-limit-tokens"`
		// RemainingRequests is the remaining requests of the rate
		// limit headers.
		RemainingRequests int `json:"x-ratelimit-remaining-requests"`
		// RemainingTokens is the remaining tokens of the rate limit
		// headers.
		RemainingTokens int `json:"x-ratelimit-remaining-tokens"`
		// ResetRequests is the reset requests of the rate limit
		// headers.
		ResetRequests ResetTime `json:"x-ratelimit-reset-requests"`
		// ResetTokens is the reset tokens of the rate limit headers.
		ResetTokens ResetTime `json:"x-ratelimit-reset-tokens"`
	}
	// ResetTime is a time.Time wrapper for the rate limit reset time.
	// string
	ResetTime string
	// Usage Represents the total token usage per request to Groq.
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
	// Endpoint is an endpoint for the groq api.
	Endpoint string

	fullURLOptions struct{ model string }
	fullURLOption  func(*fullURLOptions)
	response       interface{ SetHeader(http.Header) }
)

const (
	// FormatText is the text format. It is the default format of a
	// response.
	FormatText Format = "text"
	// FormatJSON is the JSON format. There is no support for streaming with
	// JSON format selected.
	FormatJSON Format = "json"
	// FormatSRT is the SRT format. This is a text format that is only
	// supported for the transcription API.
	// SRT format selected.
	FormatSRT Format = "srt"
	// FormatVTT is the VTT format. This is a text format that is only
	// supported for the transcription API.
	FormatVTT Format = "vtt"
	// FormatVerboseJSON is the verbose JSON format. This is a JSON format
	// that is only supported for the transcription API.
	FormatVerboseJSON Format = "verbose_json"
	// FormatJSONObject is the json object chat
	// completion response format type.
	FormatJSONObject Format = "json_object"
	// FormatJSONSchema is the json schema chat
	// completion response format type.
	FormatJSONSchema Format = "json_schema"

	// groqAPIURLv1 is the base URL for the Groq API.
	groqAPIURLv1 = "https://api.groq.com/openai/v1"

	chatCompletionsSuffix Endpoint = "/chat/completions"
	transcriptionsSuffix  Endpoint = "/audio/transcriptions"
	translationsSuffix    Endpoint = "/audio/translations"
	embeddingsSuffix      Endpoint = "/embeddings"
	moderationsSuffix     Endpoint = "/moderations"
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
	}
	for _, opt := range opts {
		opt(c)
	}
	c.header.SetCommonHeaders = func(req *http.Request) {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.groqAPIKey))
		if c.orgID != "" {
			req.Header.Set("OpenAI-Organization", c.orgID)
		}
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

func sendRequestStream[T streams.Streamer[ChatCompletionStreamResponse]](
	client *Client,
	req *http.Request,
) (*streams.StreamReader[*ChatCompletionStreamResponse], error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.client.Do(
		req,
	) //nolint:bodyclose // body is closed in stream.Close()
	if err != nil {
		return new(streams.StreamReader[*ChatCompletionStreamResponse]), err
	}
	if isFailureStatusCode(resp) {
		return new(streams.StreamReader[*ChatCompletionStreamResponse]), client.handleErrorResp(resp)
	}
	return streams.NewStreamReader[ChatCompletionStreamResponse](
		resp.Body,
		resp.Header,
		client.emptyMessagesLimit,
	), nil
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

func withModel[
	T models.ChatModel | models.AudioModel | models.ModerationModel,
](model T) fullURLOption {
	return func(args *fullURLOptions) {
		args.model = string(model)
	}
}

func (c *Client) handleErrorResp(resp *http.Response) error {
	var errRes groqerr.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errRes)
	if err != nil || errRes.Error == nil {
		reqErr := &groqerr.ErrRequest{
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
