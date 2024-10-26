package composio

import (
	"context"
	"os"
	"testing"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	if !test.IsUnitTest() {
		t.Skip()
	}
	a := assert.New(t)
	ctx := context.Background()
	key, err := test.GetAPIKey("COMPOSIO_API_KEY")
	a.NoError(err)
	client, err := NewComposer(
		key,
		WithLogger(test.DefaultLogger),
	)
	a.NoError(err)
	ts, err := client.GetTools(
		ctx, WithApp("GITHUB"), WithUseCase("StarRepo"))
	a.NoError(err)
	a.NotEmpty(ts)
	groqClient, err := groq.NewClient(
		os.Getenv("GROQ_KEY"),
	)
	a.NoError(err, "NewClient error")
	response, err := groqClient.CreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model: groq.ModelLlama3Groq8B8192ToolUsePreview,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.ChatMessageRoleUser,
				Content: "Star the facebookresearch/spiritlm repository on GitHub",
			},
		},
		MaxTokens: 2000,
		Tools:     ts,
	})
	a.NoError(err)
	a.NotEmpty(response.Choices[0].Message.ToolCalls)
	resp2, err := client.Run(ctx, response)
	a.NoError(err)
	a.NotEmpty(resp2)
	t.Logf("%+v\n", resp2)
}
