package groq

import (
	"context"
	"fmt"
	"log/slog"
)

type (
	// Agent is an agent.
	Agent struct {
		client  *Client
		logger  *slog.Logger
		runners []ToolRunner
	}
	// ToolRunner is an interface for a tool manager.
	ToolRunner interface {
		Run(
			ctx context.Context,
			response ChatCompletionResponse,
		) ([]ChatCompletionMessage, error)
	}
)

// Run runs the agent on a chat completion response.
func (a *Agent) Run(
	ctx context.Context,
	response ChatCompletionResponse,
) ([]ChatCompletionMessage, error) {
	for _, runner := range a.runners {
		messages, err := runner.Run(ctx, response)
		if err == nil {
			return messages, nil
		}
	}
	return nil, fmt.Errorf("no runners found for response")
}

// NewAgent creates a new agent.
func NewAgent(
	client *Client,
	logger *slog.Logger,
	runners ...ToolRunner,
) *Agent {
	return &Agent{
		client:  client,
		logger:  logger,
		runners: runners,
	}
}
