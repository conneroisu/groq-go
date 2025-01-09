package test

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
)

// CreateTestFile creates a fake file with "hello" as the content.
func CreateTestFile(t *testing.T, path string) {
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create file %v", err)
	}

	if _, err = file.WriteString("hello"); err != nil {
		t.Fatalf("failed to write to file %v", err)
	}
	file.Close()
}

// CreateTestDirectory creates a temporary folder which will be deleted when cleanup is called.
func CreateTestDirectory(t *testing.T) (path string, cleanup func()) {
	t.Helper()

	path, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		t.Fatalf("failed to create directory %v", err)
	}

	return path, func() { os.RemoveAll(path) }
}

// TokenRoundTripper is a struct that implements the RoundTripper
// interface, specifically to handle the authentication token by adding a token
// to the request header. We need this because the API requires that each
// request include a valid API token in the headers for authentication and
// authorization.
type TokenRoundTripper struct {
	Token    string
	Fallback http.RoundTripper
}

// RoundTrip takes an *http.Request as input and returns an
// *http.Response and an error.
//
// It is expected to use the provided request to create a connection to an HTTP
// server and return the response, or an error if one occurred. The returned
// Response should have its Body closed. If the RoundTrip method returns an
// error, the Client's Get, Head, Post, and PostForm methods return the same
// error.
func (t *TokenRoundTripper) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.Token)
	return t.Fallback.RoundTrip(req)
}

// IsIntegrationTest returns true if the unit test environment variable is set.
func IsIntegrationTest() bool {
	return os.Getenv("UNIT") != ""
}

// GetAPIKey returns the api key.
func GetAPIKey(key string) (string, error) {
	apiKey := os.Getenv(key)
	if apiKey == "" {
		return "", fmt.Errorf("api key: %s is required", key)
	}
	return apiKey, nil
}

// DefaultLogger is a default logger.
var DefaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: true,
	Level:     slog.LevelDebug,
	ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == "time" {
			return slog.Attr{}
		}
		if a.Key == "level" {
			return slog.Attr{}
		}
		if a.Key == slog.SourceKey {
			str := a.Value.String()
			split := strings.Split(str, "/")
			if len(split) > 2 {
				a.Value = slog.StringValue(strings.Join(split[len(split)-2:], "/"))
				a.Value = slog.StringValue(strings.Replace(a.Value.String(), "}", "", -1))
			}
		}
		return a
	}}))
