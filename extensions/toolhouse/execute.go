package toolhouse

import (
	"context"
	"fmt"
	"net/http"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/conneroisu/groq-go/pkg/tools"
)

type (
	request struct {
		Content  tools.ToolCall `json:"content,omitempty"`
		Provider string         `json:"provider"`
		Metadata map[string]any `json:"metadata"`
		Bundle   string         `json:"bundle"`
	}
)

// Run runs the extension on the given history.
func (e *Toolhouse) Run(
	ctx context.Context,
	response groq.ChatCompletionResponse,
) ([]groq.ChatCompletionMessage, error) {
	var respH []groq.ChatCompletionMessage
	var toolCall tools.ToolCall
	e.logger.Debug("Running Toolhouse extension", "response", response)
	if response.Choices[0].FinishReason != groq.ReasonFunctionCall && response.Choices[0].FinishReason != "tool_calls" {
		return nil, fmt.Errorf("not a function call")
	}
	for _, toolCall = range response.Choices[0].Message.ToolCalls {
		req, err := builders.NewRequest(
			ctx,
			e.header,
			http.MethodPost,
			fmt.Sprintf("%s%s", e.baseURL, runToolEndpoint),
			builders.WithBody(request{
				Content:  toolCall,
				Provider: e.provider,
				Metadata: e.metadata,
				Bundle:   e.bundle,
			}),
		)
		if err != nil {
			return nil, err
		}
		e.logger.Debug("toolhouse running tool", "tool", toolCall.Function.Name, "call", toolCall.Function.Arguments)
		var runResp struct {
			Provider string `json:"provider"`
			Content  struct {
				Role       string `json:"role"`
				ToolCallID string `json:"tool_call_id"`
				Name       string `json:"name"`
				Content    string `json:"content"`
			} `json:"content"`
		}
		err = e.sendRequest(req, &runResp)
		if err != nil {
			return nil, err
		}
		respH = append(respH, groq.ChatCompletionMessage{
			Content: runResp.Content.Content,
			Name:    runResp.Content.Name,
			Role:    groq.RoleFunction,
		})
	}
	return respH, nil
}
