package groq_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/conneroisu/groq-go"
	"github.com/stretchr/testify/assert"
)

// TestCompletionsWrongModel tests the CreateCompletion method with a wrong model.
func TestCompletionsWrongModel(t *testing.T) {
	a := assert.New(t)
	client, err := groq.NewClient(
		"whatever",
		groq.WithBaseURL("http://localhost/v1"),
	)
	a.NoError(err, "NewClient error")

	_, err = client.CreateCompletion(
		context.Background(),
		groq.CompletionRequest{
			MaxTokens: 5,
			Model:     groq.GPT3Dot5Turbo,
		},
	)
	if !errors.Is(err, groq.ErrCompletionUnsupportedModel{}) {
		t.Fatalf(
			"CreateCompletion should return ErrCompletionUnsupportedModel, but returned: %v",
			err,
		)
	}
}

// TestCompletionWithStream tests the CreateCompletion method with a stream.
func TestCompletionWithStream(t *testing.T) {
	a := assert.New(t)
	client, err := groq.NewClient(
		"whatever",
		groq.WithBaseURL("http://localhost/v1"),
	)
	a.NoError(err, "NewClient error")

	ctx := context.Background()
	req := groq.CompletionRequest{Stream: true}
	_, err = client.CreateCompletion(ctx, req)
	if !errors.Is(err, groq.ErrCompletionStreamNotSupported{}) {
		t.Fatalf(
			"CreateCompletion didn't return ErrCompletionStreamNotSupported",
		)
	}
}

// TestCompletions Tests the completions endpoint of the API using the mocked server.
func TestCompletions(t *testing.T) {
	a := assert.New(t)
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler("/v1/completions", handleCompletionEndpoint)
	req := groq.CompletionRequest{
		MaxTokens: 5,
		Model:     "ada",
		Prompt:    "Lorem ipsum",
	}
	_, err := client.CreateCompletion(context.Background(), req)
	a.NoError(err, "CreateCompletion error")
}

// handleCompletionEndpoint Handles the completion endpoint by the test server.
func handleCompletionEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error
	var resBytes []byte

	// completions only accepts POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	var completionReq groq.CompletionRequest
	if completionReq, err = getCompletionBody(r); err != nil {
		http.Error(w, "could not read request", http.StatusInternalServerError)
		return
	}
	res := groq.CompletionResponse{
		ID:      strconv.Itoa(int(time.Now().Unix())),
		Object:  "test-object",
		Created: time.Now().Unix(),
		// would be nice to validate Model during testing, but
		// this may not be possible with how much upkeep
		// would be required / wouldn't make much sense
		Model: completionReq.Model,
	}
	// create completions
	n := completionReq.N
	if n == 0 {
		n = 1
	}
	for i := 0; i < n; i++ {
		// generate a random string of length completionReq.Length
		completionStr := strings.Repeat("a", completionReq.MaxTokens)
		if completionReq.Echo {
			completionStr = completionReq.Prompt.(string) + completionStr
		}
		res.Choices = append(res.Choices, groq.CompletionChoice{
			Text:  completionStr,
			Index: i,
		})
	}
	inputTokens := numTokens(completionReq.Prompt.(string)) * n
	completionTokens := completionReq.MaxTokens * n
	res.Usage = groq.Usage{
		PromptTokens:     inputTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      inputTokens + completionTokens,
	}
	resBytes, _ = json.Marshal(res)
	fmt.Fprintln(w, string(resBytes))
}

// numTokens Returns the number of GPT-3 encoded tokens in the given text.
// This function approximates based on the rule of thumb stated by OpenAI:
// https://beta.openai.com/tokenizer
//
// TODO: implement an actual tokenizer for each model available and use that
// instead.
func numTokens(s string) int {
	return int(float32(len(s)) / 4)
}

// getCompletionBody Returns the body of the request to create a completion.
func getCompletionBody(r *http.Request) (groq.CompletionRequest, error) {
	completion := groq.CompletionRequest{}
	// read the request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return groq.CompletionRequest{}, err
	}
	err = json.Unmarshal(reqBody, &completion)
	if err != nil {
		return groq.CompletionRequest{}, err
	}
	return completion, nil
}
