package groq

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/conneroisu/groq-go/pkg/groqerr"
	"github.com/conneroisu/groq-go/internal/streams"
)

//go:generate go run ./cmd/generate-models
//go:generate go run github.com/princjef/gomarkdoc/cmd/gomarkdoc@v1.1.0 -o README.md -e .

type (
	// Client is a Groq api client.
	Client struct {
		// Groq API key
		groqAPIKey         string
		orgID              string
		baseURL            string
		emptyMessagesLimit uint

		header             builders.Header
		requestFormBuilder builders.FormBuilder

		client *http.Client
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
func (c *Client) fullURL(suffix endpoint, setters ...fullURLOption) string {
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
	T ChatModel | AudioModel | ModerationModel,
](model T) fullURLOption {
	return func(args *fullURLOptions) {
		args.model = string(model)
	}
}
