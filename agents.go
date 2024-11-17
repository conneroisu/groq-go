package groq

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/conneroisu/groq-go/pkg/tools"
)

type (
	// Agency is a collection of agents.
	Agency struct {
		client *Client
		agents []Agent
		logger *slog.Logger
	}
	// Agent is an agent.
	Agent struct {
		client    *Client
		logger    *slog.Logger
		providers []ToolProvider
		history   []ChatCompletionMessage
		inbox     chan ChatCompletionMessage
	}
	// ToolProvider is an interface for a tool provider.
	ToolProvider interface {
		// Run runs responded tool calls.
		Run(
			ctx context.Context,
			response ChatCompletionResponse,
		) ([]ChatCompletionMessage, error)
		// Get gets the tools for the provider.
		Get(
			ctx context.Context,
		) ([]tools.Tool, error)
		// Resolve resolves a tool call.
		//
		// Implementations must not return an error if the tool is not found.
		Resolve(
			ctx context.Context,
			call tools.ToolCall,
		) (bool, error)
	}
)

// NewAgent creates a new agent.
func NewAgent(
	client *Client,
	logger *slog.Logger,
	tools ...ToolProvider,
) *Agent {
	return &Agent{
		client:    client,
		logger:    logger,
		providers: tools,
	}
}

func (agent *Agent) resolveTool(
	ctx context.Context,
	call tools.ToolCall,
) (ToolProvider, error) {
	for _, provider := range agent.providers {
		ok, err := provider.Resolve(ctx, call)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		return provider, nil
	}
	return nil, fmt.Errorf("tool not found")
}

// run runs the agent on a chat completion response.
//
// More specifically, it runs the tool calls made by the model.
func (agent *Agent) run(
	ctx context.Context,
	response ChatCompletionResponse,
) (hist []ChatCompletionMessage, err error) {
	for _, tool := range response.Choices[0].Message.ToolCalls {
		runner, err := agent.resolveTool(ctx, tool)
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
func (agent *Agent) refresh(
	ctx context.Context,
) error {
	// TODO: refresh agent context
	return nil
}

// Start starts the agency within the given context.
func (agency *Agency) Start(ctx context.Context) error {
	// TODO: start agents
	return nil
}
