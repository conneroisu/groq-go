package groq_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	groq "github.com/conneroisu/groq-go"
	"github.com/stretchr/testify/assert"
)

// TestModerate tests the Moderate method of the client.
func TestModerate(t *testing.T) {
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler(
		"/v1/chat/completions",
		handleModerationEndpoint,
	)
	mod, err := client.Moderate(context.Background(), groq.ModerationRequest{
		Model: groq.ModerationTextStable,
		Input: "I want to kill them.",
	})
	a := assert.New(t)
	a.NoError(err, "Moderation error")
	a.Equal(true, mod.Flagged)
	a.Equal(
		mod.Categories,
		[]groq.HarmfulCategory{
			groq.CategoryViolentCrimes,
			groq.CategoryNonviolentCrimes,
		},
	)
}

// handleModerationEndpoint handles the moderation endpoint.
func handleModerationEndpoint(w http.ResponseWriter, r *http.Request) {
	response := groq.ChatCompletionResponse{
		ID:      "chatcmpl-123",
		Object:  "chat.completion",
		Created: 1693721698,
		Model:   groq.ModerationTextStable,
		Choices: []groq.ChatCompletionChoice{
			{
				Message: groq.ChatCompletionMessage{
					Role:    groq.ChatMessageRoleAssistant,
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
