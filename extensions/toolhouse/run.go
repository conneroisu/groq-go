package toolhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/builders"
)

type (
	request struct {
		Content  groq.ToolCall  `json:"content,omitempty"`
		Provider string         `json:"provider"`
		Metadata map[string]any `json:"metadata"`
		Bundle   string         `json:"bundle"`
	}
)

// MustRun runs the extension on the given history.
//
// It panics if an error occurs.
func (e *Extension) MustRun(
	ctx context.Context,
	response groq.ChatCompletionResponse,
) []groq.ChatCompletionMessage {
	respH, err := e.Run(ctx, response)
	if err != nil {
		panic(err)
	}
	return respH
}

// Run runs the extension on the given history.
func (e *Extension) Run(
	ctx context.Context,
	response groq.ChatCompletionResponse,
) ([]groq.ChatCompletionMessage, error) {
	e.logger.Debug("Running Toolhouse extension", "response", response)
	if response.Choices[0].FinishReason != groq.FinishReasonFunctionCall && response.Choices[0].FinishReason != "tool_calls" {
		return nil, fmt.Errorf("Not a function call")
	}
	respH := []groq.ChatCompletionMessage{}
	for _, tool := range response.Choices[0].Message.ToolCalls {
		req, err := builders.NewRequest(
			ctx,
			e.header,
			http.MethodPost,
			fmt.Sprintf("%s%s", e.baseURL, runToolEndpoint),
			builders.WithBody(request{
				Content:  tool,
				Provider: e.provider,
				Metadata: e.metadata,
				Bundle:   e.bundle,
			}),
		)
		if err != nil {
			return nil, err
		}
		resp, err := e.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%v", resp)
		}
		bdy, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var runResp struct {
			Provider string `json:"provider"`
			Content  struct {
				Role       string `json:"role"`
				ToolCallID string `json:"tool_call_id"`
				Name       string `json:"name"`
				Content    string `json:"content"`
			} `json:"content"`
		}
		err = json.Unmarshal(bdy, &runResp)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %w: %s", err, string(bdy))
		}
		respH = append(respH, groq.ChatCompletionMessage{
			Content: runResp.Content.Content,
			Name:    runResp.Content.Name,
			Role:    groq.ChatMessageRoleFunction,
		})
	}
	return respH, nil
}
