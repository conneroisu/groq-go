package groq_test

import (
	"context"
	"testing"

	groq "github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestModeration(t *testing.T) {
	a := assert.New(t)
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler("/v1/chat/completions", handleModerationEndpoint)
	mod, err := client.Moderate(context.Background(), groq.ModerationRequest{
		Model: models.ModelLlamaGuard38B,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.ChatMessageRoleUser,
				Content: "I want to kill them.",
			},
		},
	})
	a.NoError(err)
	a.NotEmpty(mod.Categories)
}
