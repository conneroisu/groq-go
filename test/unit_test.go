package test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/conneroisu/groq-go"
	"github.com/stretchr/testify/assert"
)

func TestTestServer(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	client, err := groq.NewClient(os.Getenv("GROQ_KEY"))
	a.NoError(err, "NewClient error")
	strm, err := client.CreateChatCompletionStream(ctx, groq.ChatCompletionRequest{
		Model: groq.Llama3070B8192ToolUsePreview,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.ChatMessageRoleUser,
				Content: "Hello! What is a proface industrial display?",
			},
		},
		MaxTokens: 90,
		Stream:    true,
	})
	a.NoError(err, "CreateCompletionStream error")

	for {
		val, err := strm.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		t.Logf("%s\n", val.Choices[0].Delta.Content)
	}
}
