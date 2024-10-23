package composio

import (
	"context"
	"fmt"

	"github.com/conneroisu/groq-go"
)

// Run runs the composio client on a chat completion response.
func (c *Composio) Run(
	ctx context.Context,
	response groq.ChatCompletionResponse,
) error {
	if response.Choices[0].FinishReason != groq.FinishReasonFunctionCall && response.Choices[0].FinishReason != "tool_calls" {
		return fmt.Errorf("not a function call")
	}
	// for toolCall := range response.Choices[0].ToolCalls {
	//         callURL := fmt.Sprintf("%s/%s/execute", c.baseURL, toolCall.ID)
	// }
	return nil
}
