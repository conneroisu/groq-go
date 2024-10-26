package composio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/builders"
)

type (
	// Runner is an interface for composio run.
	Runner interface {
		Run(ctx context.Context, response groq.ChatCompletionResponse) (
			[]groq.ChatCompletionMessage, error)
	}
	request struct {
		ConnectedAccountID string         `json:"connectedAccountId"`
		EntityID           string         `json:"entityId"`
		AppName            string         `json:"appName"`
		Input              map[string]any `json:"input"`
		Text               string         `json:"text,omitempty"`
		AuthConfig         map[string]any `json:"authConfig,omitempty"`
	}
)

// Run runs the composio client on a chat completion response.
func (c *Composio) Run(
	ctx context.Context,
	response groq.ChatCompletionResponse,
) ([]groq.ChatCompletionMessage, error) {
	var respH []groq.ChatCompletionMessage
	if response.Choices[0].FinishReason != groq.FinishReasonFunctionCall &&
		response.Choices[0].FinishReason != "tool_calls" {
		return nil, fmt.Errorf("Not a function call")
	}
	connectedAccount, err := c.GetConnectedAccounts(ctx, WithShowActiveOnly(true))
	if err != nil {
		return nil, err
	}
	c.logger.Debug("connected accounts", "accounts", connectedAccount)
	for _, toolCall := range response.Choices[0].Message.ToolCalls {
		var args map[string]any
		if json.Valid([]byte(toolCall.Function.Arguments)) {
			err = json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
			if err != nil {
				return nil, err
			}
			c.logger.Debug("arguments", "args", args)
		}
		req, err := builders.NewRequest(
			ctx,
			c.header,
			http.MethodPost,
			fmt.Sprintf("%s/v2/actions/%s/execute", c.baseURL, toolCall.Function.Name),
			builders.WithBody(&request{
				ConnectedAccountID: connectedAccount[0].ID,
				EntityID:           "default",
				AppName:            toolCall.Function.Name,
				Input:              args,
				Text:               "",
				AuthConfig:         map[string]any{},
			}),
		)
		if err != nil {
			return nil, err
		}
		var body string
		err = c.doRequest(req, &body)
		if err != nil {
			return nil, err
		}
		respH = append(respH, groq.ChatCompletionMessage{
			Content: string(body),
			Name:    toolCall.ID,
			Role:    groq.ChatMessageRoleFunction,
		})
	}
	return respH, nil
}
