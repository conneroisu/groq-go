package toolhouse_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/extensions/toolhouse"
	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/conneroisu/groq-go/pkg/tools"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	ts := test.NewTestServer()
	ts.RegisterHandler("/run_tools", func(w http.ResponseWriter, r *http.Request) {
		var runResp struct {
			Provider string `json:"provider"`
			Content  struct {
				Role       string `json:"role"`
				ToolCallID string `json:"tool_call_id"`
				Name       string `json:"name"`
				Content    string `json:"content"`
			} `json:"content"`
		}
		runResp.Content.Content = "response1"
		runResp.Content.Name = "tool"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonBytes, err := json.Marshal(runResp)
		a.NoError(err)
		_, err = w.Write(jsonBytes)
		a.NoError(err)
	})
	testS := ts.ToolhouseTestServer()
	testS.Start()
	client, err := toolhouse.NewExtension(
		test.GetTestToken(),
		toolhouse.WithBaseURL(testS.URL),
		toolhouse.WithClient(testS.Client()),
		toolhouse.WithLogger(test.DefaultLogger),
		toolhouse.WithMetadata(map[string]any{
			"id":       "conner",
			"timezone": 5,
		}),
	)
	a.NoError(err)
	history := []groq.ChatCompletionMessage{
		{
			Role:    groq.RoleUser,
			Content: "",
			ToolCalls: []tools.ToolCall{
				{
					Function: tools.FunctionCall{
						Name: "tool",
					},
				},
			},
		},
	}
	resp, err := client.Run(ctx, groq.ChatCompletionResponse{
		Choices: []groq.ChatCompletionChoice{
			{
				Message:      history[0],
				FinishReason: groq.ReasonFunctionCall,
			},
		},
	})
	a.NoError(err)
	assert.Equal(t, "response1", resp[0].Content)
}
