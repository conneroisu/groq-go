package test_test

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

// TestCreateTestFile verifies that CreateTestFile correctly creates a file with the expected content.
func TestCreateTestFile(t *testing.T) {
	a := assert.New(t)

	// Create a temporary directory and ensure cleanup.
	dir, cleanup := test.CreateTestDirectory(t)
	defer cleanup()

	// Define the path for the test file.
	filePath := filepath.Join(dir, "testfile.txt")

	// Call the function under test.
	test.CreateTestFile(t, filePath)

	// Check that the file exists.
	info, err := os.Stat(filePath)
	a.NoError(err, "File should exist")
	a.False(info.IsDir(), "Should be a file, not a directory")

	// Read and verify the file content.
	content, err := os.ReadFile(filePath)
	a.NoError(err, "Should be able to read the file")
	a.Equal("hello", string(content), "File content should be 'hello'")
}

// TestCreateTestDirectory ensures that CreateTestDirectory creates a directory and the cleanup function removes it.
func TestCreateTestDirectory(t *testing.T) {
	a := assert.New(t)

	// Create the test directory.
	dir, cleanup := test.CreateTestDirectory(t)

	// Check that the directory exists.
	info, err := os.Stat(dir)
	a.NoError(err, "Directory should exist")
	a.True(info.IsDir(), "Should be a directory")

	// Write a test file inside the directory.
	testFilePath := filepath.Join(dir, "test.txt")
	err = os.WriteFile(testFilePath, []byte("test content"), 0644)
	a.NoError(err, "Should be able to write a file in the directory")

	// Perform cleanup.
	cleanup()

	// Verify that the directory has been removed.
	_, err = os.Stat(dir)
	a.True(os.IsNotExist(err), "Directory should be deleted after cleanup")
}

// MockRoundTripper is a mock implementation of http.RoundTripper for testing purposes.
type MockRoundTripper struct {
	LastRequest *http.Request
}

// RoundTrip captures the request and returns a dummy response.
func (m *MockRoundTripper) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	m.LastRequest = req
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("OK")),
		Header:     make(http.Header),
	}, nil
}

// TestTokenRoundTripper verifies that TokenRoundTripper adds the correct Authorization header.
func TestTokenRoundTripper(t *testing.T) {
	a := assert.New(t)

	// Prepare the mock fallback RoundTripper.
	mockRT := &MockRoundTripper{}

	// Initialize the TokenRoundTripper with a test token.
	tokenRT := &test.TokenRoundTripper{
		Token:    "test-token",
		Fallback: mockRT,
	}

	// Create an HTTP client using the TokenRoundTripper.
	client := &http.Client{
		Transport: tokenRT,
	}

	// Prepare a test HTTP request.
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	a.NoError(err, "Should be able to create a new request")

	// Perform the HTTP request.
	resp, err := client.Do(req)
	a.NoError(err, "HTTP request should succeed")
	a.Equal(
		http.StatusOK,
		resp.StatusCode,
		"Response status code should be 200",
	)

	// Verify that the Authorization header was added.
	a.NotNil(mockRT.LastRequest, "LastRequest should be captured")
	authHeader := mockRT.LastRequest.Header.Get("Authorization")
	a.Equal(
		"Bearer test-token",
		authHeader,
		"Authorization header should contain the correct token",
	)
}
