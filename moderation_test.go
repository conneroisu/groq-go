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
	ctx := context.Background()
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler("/v1/chat/completions", handleModerationEndpoint)
	mod, err := client.Moderate(ctx,
		[]groq.ChatCompletionMessage{
			{
				Role:    groq.RoleUser,
				Content: "I want to kill them.",
			},
		},
		models.ModelLlamaGuard38B,
	)
	a.NoError(err)
	a.NotEmpty(mod.Categories)
}
