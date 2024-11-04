package groq

import (
	"context"
	"log/slog"

	"github.com/conneroisu/groq-go/pkg/tools"
)

type (
	// Agenter is an interface for an agent.
	Agenter interface {
		ToolManager
	}
	// ToolManager is an interface for a tool manager.
	ToolManager interface {
		ToolGetter
		ToolRunner
	}
	// ToolGetter is an interface for a tool getter.
	ToolGetter interface {
		Get(
			ctx context.Context,
			params ToolGetParams,
		) ([]tools.Tool, error)
	}
	// ToolRunner is an interface for a tool runner.
	ToolRunner interface {
		Run(
			ctx context.Context,
			response ChatCompletionResponse,
		) ([]ChatCompletionMessage, error)
	}
	// ToolGetParams are the parameters for getting tools.
	ToolGetParams struct {
	}
	// Router is an agent router.
	//
	// It is used to route messages to the appropriate model.
	Router struct {
		// Agents is the agents of the router.
		Agents []Agent
		// Logger is the logger of the router.
		Logger *slog.Logger
	}
)

// Agent is an agent.
type Agent struct {
	client *Client
	logger *slog.Logger
}

// NewAgent creates a new agent.
func NewAgent(
	client *Client,
	logger *slog.Logger,
) *Agent {
	return &Agent{
		client: client,
		logger: logger,
	}
}
