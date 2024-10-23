package composio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/builders"
)

// Run runs the composio client on a chat completion response.
func (c *Composio) Run(
	ctx context.Context,
	response groq.ChatCompletionResponse,
) ([]groq.ChatCompletionMessage, error) {
	var respH []groq.ChatCompletionMessage
	var bdy []byte
	if response.Choices[0].FinishReason != groq.FinishReasonFunctionCall && response.Choices[0].FinishReason != "tool_calls" {
		return nil, fmt.Errorf("Not a function call")
	}
	for _, toolCall := range response.Choices[0].Message.ToolCalls {
		callURL := fmt.Sprintf("%s/%s/execute", c.baseURL, toolCall.ID)
		req, err := builders.NewRequest(
			ctx,
			c.header,
			http.MethodPost,
			callURL,
			builders.WithBody(toolCall.Function.Arguments),
		)
		if err != nil {
			return nil, err
		}
		var toolResp struct {
			Properties struct {
				Data       interface{} `json:"data"`
				Successful interface{} `json:"successful"`
				Error      interface{} `json:"error"`
			} `json:"properties"`
		}
		err = c.doRequest(req, &toolResp)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bdy, &toolResp)
		if err != nil {
			return nil, err
		}
		respH = append(respH, groq.ChatCompletionMessage{
			Content: string(bdy),
			Name:    toolCall.ID,
			Role:    groq.ChatMessageRoleFunction,
		})
	}
	return respH, nil
}
