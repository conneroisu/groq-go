// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package groq_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/internal/testutil"
	"github.com/conneroisu/groq-go/option"
	"github.com/conneroisu/groq-go/shared"
)

func TestChatCompletionNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Chat.Completions.New(context.TODO(), groq.ChatCompletionNewParams{
		Messages: groq.F([]groq.ChatCompletionNewParamsMessageUnion{groq.ChatCompletionNewParamsMessagesChatCompletionRequestSystemMessage{
			Content: groq.F("content"),
			Role:    groq.F(groq.ChatCompletionNewParamsMessagesChatCompletionRequestSystemMessageRoleSystem),
			Name:    groq.F("name"),
		}}),
		Model:            groq.F(groq.ChatCompletionNewParamsModelGemma7bIt),
		FrequencyPenalty: groq.F(-2.000000),
		FunctionCall:     groq.F[groq.ChatCompletionNewParamsFunctionCallUnion](groq.ChatCompletionNewParamsFunctionCallString(groq.ChatCompletionNewParamsFunctionCallStringNone)),
		Functions: groq.F([]groq.ChatCompletionNewParamsFunction{{
			Name:        groq.F("name"),
			Description: groq.F("description"),
			Parameters: groq.F(map[string]interface{}{
				"foo": "bar",
			}),
		}, {
			Name:        groq.F("name"),
			Description: groq.F("description"),
			Parameters: groq.F(map[string]interface{}{
				"foo": "bar",
			}),
		}, {
			Name:        groq.F("name"),
			Description: groq.F("description"),
			Parameters: groq.F(map[string]interface{}{
				"foo": "bar",
			}),
		}}),
		LogitBias: groq.F(map[string]int64{
			"foo": int64(0),
		}),
		Logprobs:          groq.F(true),
		MaxTokens:         groq.F(int64(0)),
		N:                 groq.F(int64(1)),
		ParallelToolCalls: groq.F(true),
		PresencePenalty:   groq.F(-2.000000),
		ResponseFormat: groq.F(groq.ChatCompletionNewParamsResponseFormat{
			Type: groq.F(groq.ChatCompletionNewParamsResponseFormatTypeText),
		}),
		Seed:        groq.F(int64(0)),
		Stop:        groq.F[groq.ChatCompletionNewParamsStopUnion](shared.UnionString("\n")),
		Stream:      groq.F(true),
		Temperature: groq.F(1.000000),
		ToolChoice:  groq.F[groq.ChatCompletionNewParamsToolChoiceUnion](groq.ChatCompletionNewParamsToolChoiceString(groq.ChatCompletionNewParamsToolChoiceStringNone)),
		Tools: groq.F([]groq.ChatCompletionNewParamsTool{{
			Function: groq.F(groq.ChatCompletionNewParamsToolsFunction{
				Name:        groq.F("name"),
				Description: groq.F("description"),
				Parameters: groq.F(map[string]interface{}{
					"foo": "bar",
				}),
			}),
			Type: groq.F(groq.ChatCompletionNewParamsToolsTypeFunction),
		}, {
			Function: groq.F(groq.ChatCompletionNewParamsToolsFunction{
				Name:        groq.F("name"),
				Description: groq.F("description"),
				Parameters: groq.F(map[string]interface{}{
					"foo": "bar",
				}),
			}),
			Type: groq.F(groq.ChatCompletionNewParamsToolsTypeFunction),
		}, {
			Function: groq.F(groq.ChatCompletionNewParamsToolsFunction{
				Name:        groq.F("name"),
				Description: groq.F("description"),
				Parameters: groq.F(map[string]interface{}{
					"foo": "bar",
				}),
			}),
			Type: groq.F(groq.ChatCompletionNewParamsToolsTypeFunction),
		}}),
		TopLogprobs: groq.F(int64(0)),
		TopP:        groq.F(1.000000),
		User:        groq.F("user"),
	})
	if err != nil {
		var apierr *groq.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
