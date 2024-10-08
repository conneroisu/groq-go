package toolhouse

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpRequestBuilder_Build_WithBody(t *testing.T) {
	// Setup
	ctx := context.Background()
	url := "https://example.com"
	builder := newRequestBuilder()

	bodyContent := map[string]string{
		"key": "value",
	}

	// Test building a request with body
	req, err := builder.Build(ctx, http.MethodPost, url, bodyContent, nil)
	assert.NoError(t, err, "Expected no error while building request")
	assert.Equal(t, http.MethodPost, req.Method, "Expected POST method")
	assert.Equal(t, url, req.URL.String(), "Expected correct URL")

	// Check the body
	expectedBody, err := json.Marshal(bodyContent)
	assert.NoError(t, err, "Expected no error while marshalling body")
	bodyBytes, err := io.ReadAll(req.Body)
	assert.NoError(t, err, "Expected no error reading request body")
	assert.Equal(t, expectedBody, bodyBytes, "Expected body to match")
}

func TestHttpRequestBuilder_Build_WithoutBody(t *testing.T) {
	// Setup
	ctx := context.Background()
	url := "https://example.com"
	builder := newRequestBuilder()

	// Test building a request without body
	req, err := builder.Build(ctx, http.MethodGet, url, nil, nil)
	assert.NoError(t, err, "Expected no error while building request")
	assert.Equal(t, http.MethodGet, req.Method, "Expected GET method")
	assert.Equal(t, url, req.URL.String(), "Expected correct URL")

	// Check that there is no body
	assert.Nil(t, req.Body, "Expected no body for GET request")
}

func TestHttpRequestBuilder_Build_WithHeaders(t *testing.T) {
	// Setup
	ctx := context.Background()
	url := "https://example.com"
	builder := newRequestBuilder()

	headers := http.Header{}
	headers.Set("X-Custom-Header", "CustomValue")

	// Test building a request with headers
	req, err := builder.Build(ctx, http.MethodGet, url, nil, headers)
	assert.NoError(t, err, "Expected no error while building request")
	assert.Equal(t, "CustomValue", req.Header.Get("X-Custom-Header"), "Expected header to be set")
}

func TestWithBody(t *testing.T) {
	// Setup
	bodyContent := "This is a test body"
	options := &requestOptions{}

	// Test withBody function
	withBody(bodyContent)(options)

	// Check that the body is set
	assert.Equal(t, bodyContent, options.body, "Expected body to be set correctly")
}

func TestHttpRequestBuilder_Build_WithReaderBody(t *testing.T) {
	// Setup
	ctx := context.Background()
	url := "https://example.com"
	builder := newRequestBuilder()

	bodyReader := bytes.NewBufferString("test body")

	// Test building a request with io.Reader as body
	req, err := builder.Build(ctx, http.MethodPost, url, bodyReader, nil)
	assert.NoError(t, err, "Expected no error while building request")
	assert.Equal(t, http.MethodPost, req.Method, "Expected POST method")
	assert.Equal(t, url, req.URL.String(), "Expected correct URL")

	// Check the body
	bodyBytes, err := io.ReadAll(req.Body)
	assert.NoError(t, err, "Expected no error reading request body")
	assert.Equal(t, "test body", string(bodyBytes), "Expected body to match the input")
}
