package groq

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/conneroisu/groq-go/pkg/tools"
)

type (
	ToolsIdx struct {
		Start int
		End   int
	}
	// Agent is an agent.
	Agent struct {
		client *Client
		logger *slog.Logger
		tools  []ToolProvider
		hist   []ChatCompletionMessage
	}
	// ToolProvider is an interface for a tool provider.
	ToolProvider interface {
		ToolRunner
		ToolGetter
		ToolResolver
	}
	// ToolRunner is an interface for a tool manager.
	ToolRunner interface {
		Run(
			ctx context.Context,
			response ChatCompletionResponse,
		) ([]ChatCompletionMessage, error)
	}
	// ToolResolver is an interface for a tool resolver.
	//
	// Implementations must not return an error if the tool is not found.
	ToolResolver interface {
		Resolve(
			ctx context.Context,
			call tools.ToolCall,
		) (ToolRunner, error)
	}
	// ToolGetter is an interface for a tool getter.
	ToolGetter interface {
		Get(
			ctx context.Context,
		) ([]tools.Tool, error)
	}
)

// NewAgent creates a new agent.
func NewAgent(
	client *Client,
	logger *slog.Logger,
	tools ...ToolProvider,
) *Agent {
	return &Agent{
		client: client,
		logger: logger,
	}
}

func (a *Agent) resolveTool(
	ctx context.Context,
	call tools.ToolCall,
) (ToolRunner, error) {
	for _, provider := range a.tools {
		runner, err := provider.Resolve(ctx, call)
		if err != nil {
			return nil, err
		}
		if runner != nil {
			continue
		}
		return runner, nil
	}
	return nil, fmt.Errorf("tool not found")
}

// Run runs the agent on a chat completion response.
//
// More specifically, it runs the tool calls made by the model.
func (a *Agent) Run(
	ctx context.Context,
	response ChatCompletionResponse,
) (hist []ChatCompletionMessage, err error) {
	for _, tool := range response.Choices[0].Message.ToolCalls {
		runner, err := a.resolveTool(ctx, tool)
		if err != nil {
			return nil, err
		}
		resp, err := runner.Run(ctx, response)
		if err != nil {
			return nil, err
		}
		hist = append(hist, resp...)
	}
	return hist, fmt.Errorf("no runners found for response")
}

// refresh refreshes the agent token context allowing longer tasks to be
// completed by the agent.
//
// It basically resets the agent's history by having a new conversation with
// the model prefixed with a summary of the previous conversation/history.
func (a *Agent) refresh(
	ctx context.Context,
) error {
	// TODO: refresh agent context
	return nil
}
