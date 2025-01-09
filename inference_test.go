package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/conneroisu/groq-go/internal/test"
	"github.com/conneroisu/groq-go/pkg/tools"
	"github.com/stretchr/testify/assert"
)

func TestChat(t *testing.T) {
	ctx := context.Background()
	a := assert.New(t)
	ts := test.NewTestServer()
	returnObj := ChatCompletionResponse{
		ID:      "chatcmpl-123",
		Object:  "chat.completion.chunk",
		Created: 1693721698,
		Model:   "llama3-groq-70b-8192-tool-use-preview",
		Choices: []ChatCompletionChoice{
			{
				Index: 0,
				Message: ChatCompletionMessage{
					Role:    RoleAssistant,
					Content: "Hello!",
				},
			},
		},
	}
	ts.RegisterHandler("/v1/chat/completions", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsval, err := json.Marshal(returnObj)
		a.NoError(err)
		_, err = w.Write(jsval)
		if err != nil {
			t.Fatal(err)
		}
	})
	testS := ts.GroqTestServer()
	testS.Start()
	client, err := NewClient(
		test.GetTestToken(),
		WithBaseURL(testS.URL+"/v1"),
	)
	a.NoError(err)
	resp, err := client.ChatCompletion(ctx, ChatCompletionRequest{
		Model: ModelLlama3Groq70B8192ToolUsePreview,
		Messages: []ChatCompletionMessage{
			{
				Role:    RoleUser,
				Content: "Hello!",
			},
		},
		MaxTokens: 2000,
		Tools:     []tools.Tool{},
	})
	a.NoError(err)
	a.NotEmpty(resp.Choices[0].Message.Content)
}

func TestAudioWithFailingFormBuilder(t *testing.T) {
	a := assert.New(t)
	dir, cleanup := test.CreateTestDirectory(t)
	defer cleanup()
	path := filepath.Join(dir, "fake.mp3")
	test.CreateTestFile(t, path)

	req := AudioRequest{
		FilePath:    path,
		Prompt:      "test",
		Temperature: 0.5,
		Language:    "en",
		Format:      FormatSRT,
	}

	mockFailedErr := fmt.Errorf("mock form builder fail")
	mockBuilder := &mockFormBuilder{}

	mockBuilder.mockCreateFormFile = func(string, *os.File) error {
		return mockFailedErr
	}
	err := audioMultipartForm(req, mockBuilder)
	a.ErrorIs(
		err,
		mockFailedErr,
		"audioMultipartForm should return error if form builder fails",
	)

	mockBuilder.mockCreateFormFile = func(string, *os.File) error {
		return nil
	}

	var failForField string
	mockBuilder.mockWriteField = func(fieldname, _ string) error {
		if fieldname == failForField {
			return mockFailedErr
		}
		return nil
	}

	failOn := []string{
		"model",
		"prompt",
		"temperature",
		"language",
		"response_format",
	}
	for _, failingField := range failOn {
		failForField = failingField
		mockFailedErr = fmt.Errorf(
			"mock form builder fail on field %s",
			failingField,
		)

		err = audioMultipartForm(req, mockBuilder)
		a.Error(
			err,
			mockFailedErr,
			"audioMultipartForm should return error if form builder fails",
		)
	}
}

func TestModeration(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler("/v1/chat/completions", handleModerationEndpoint)
	mod, err := client.Moderate(ctx,
		[]ChatCompletionMessage{
			{
				Role:    RoleUser,
				Content: "I want to kill them.",
			},
		},
		ModelLlamaGuard38B,
	)
	a.NoError(err)
	a.NotEmpty(mod)
}

func setupGroqTestServer() (
	client *Client,
	server *test.ServerTest,
	teardown func(),
) {
	server = test.NewTestServer()
	ts := server.GroqTestServer()
	ts.Start()
	teardown = ts.Close
	client, err := NewClient(
		test.GetTestToken(),
		WithBaseURL(ts.URL+"/v1"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// handleModerationEndpoint handles the moderation endpoint.
func handleModerationEndpoint(w http.ResponseWriter, _ *http.Request) {
	response := ChatCompletionResponse{
		ID:      "chatcmpl-123",
		Object:  "chat.completion",
		Created: 1693721698,
		Model:   ChatModel(ModelLlamaGuard38B),
		Choices: []ChatCompletionChoice{
			{
				Message: ChatCompletionMessage{
					Role:    RoleAssistant,
					Content: "unsafe\nS1,S2",
				},
				FinishReason: "stop",
			},
		},
	}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(response)
	if err != nil {
		http.Error(
			w,
			"could not encode response",
			http.StatusInternalServerError,
		)
		return
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(
			w,
			"could not write response",
			http.StatusInternalServerError,
		)
		return
	}
}

func TestCreateFileField(t *testing.T) {
	a := assert.New(t)
	t.Run("createFileField failing file", func(t *testing.T) {
		t.Parallel()
		dir, cleanup := test.CreateTestDirectory(t)
		defer cleanup()
		path := filepath.Join(dir, "fake.mp3")
		test.CreateTestFile(t, path)
		req := AudioRequest{
			FilePath: path,
		}
		mockFailedErr := fmt.Errorf("mock form builder fail")
		mockBuilder := &mockFormBuilder{
			mockCreateFormFile: func(string, *os.File) error {
				return mockFailedErr
			},
		}
		err := audioMultipartForm(req, mockBuilder)
		a.ErrorIs(
			err,
			mockFailedErr,
			"createFileField using a file should return error if form builder fails",
		)
	})

	t.Run("createFileField failing reader", func(t *testing.T) {
		t.Parallel()
		req := AudioRequest{
			FilePath: "test.wav",
			Reader:   bytes.NewBuffer([]byte(`wav test contents`)),
		}

		mockFailedErr := fmt.Errorf("mock form builder fail")
		mockBuilder := &mockFormBuilder{
			mockCreateFormFileReader: func(string, io.Reader, string) error {
				return mockFailedErr
			},
		}

		err := audioMultipartForm(req, mockBuilder)
		a.ErrorIs(
			err,
			mockFailedErr,
			"createFileField using a reader should return error if form builder fails",
		)
	})

	t.Run("createFileField failing open", func(t *testing.T) {
		t.Parallel()
		req := AudioRequest{
			FilePath: "non_existing_file.wav",
		}
		mockBuilder := builders.NewFormBuilder(&test.FailingErrorBuffer{})
		err := audioMultipartForm(req, mockBuilder)
		a.Error(
			err,
			"createFileField using file should return error when open file fails",
		)
	})
}

// mockFormBuilder is a mock form builder.
type mockFormBuilder struct {
	mockCreateFormFile       func(string, *os.File) error
	mockCreateFormFileReader func(string, io.Reader, string) error
	mockWriteField           func(string, string) error
	mockClose                func() error
}

// CreateFormFile is a mock form builder create form file method.
func (fb *mockFormBuilder) CreateFormFile(
	fieldname string,
	file *os.File,
) error {
	return fb.mockCreateFormFile(fieldname, file)
}

// CreateFormFileReader is a mock form builder create form file reader method
func (fb *mockFormBuilder) CreateFormFileReader(
	fieldname string,
	r io.Reader,
	filename string,
) error {
	return fb.mockCreateFormFileReader(fieldname, r, filename)
}

// WriteField is a mock form builder write field method.
func (fb *mockFormBuilder) WriteField(fieldname, value string) error {
	return fb.mockWriteField(fieldname, value)
}

// Close is a mock form builder close method.
func (fb *mockFormBuilder) Close() error {
	return fb.mockClose()
}

// FormDataContentType is a mock form builder.
func (fb *mockFormBuilder) FormDataContentType() string {
	return ""
}
