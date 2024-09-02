// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package groq_test

import (
	"context"
	"os"
	"testing"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/internal/testutil"
	"github.com/conneroisu/groq-go/option"
)

func TestUsage(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := groq.NewClient(
		option.WithBaseURL(baseURL),
		option.WithBearerToken("My Bearer Token"),
	)
	chatCompletionNewResponse, err := client.Chat.Completions.New(context.TODO(), groq.ChatCompletionNewParams{
		Messages: groq.F([]groq.ChatCompletionNewParamsMessageUnion{groq.ChatCompletionNewParamsMessagesChatCompletionRequestSystemMessage{
			Content: groq.F("content"),
			Role:    groq.F(groq.ChatCompletionNewParamsMessagesChatCompletionRequestSystemMessageRoleSystem),
		}}),
		Model: groq.F(groq.ChatCompletionNewParamsModelGemma7bIt),
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v\n", chatCompletionNewResponse.ID)
}
