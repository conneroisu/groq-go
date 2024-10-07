package e2b

import (
	"context"
	"io"

	"github.com/conneroisu/groq-go/pkg/mime"
)

type (
	// Kernel is a code kernel.
	//
	// It is effectively a separate runtime environment inside of a sandbox.
	//
	// You can imagine kernel as a separate environment where code is
	// executed.
	//
	// You can have multiple kernels running at the same time.
	//
	// Each kernel has its own state, so you can have multiple kernels
	// running different code at the same time.
	//
	// A kernel will be kept alive with the sandbox even if you disconnect.
	// So, it may be useful to defer the shutdown of the kernel.
	Kernel struct {
		ext *Sandbox
	}
	// ExecuteCellResponse represents the response of the execute cell method.
	ExecuteCellResponse struct {
		// CellID is the cell id of the response.
		CellID string `json:"cell_id"`
		// ExecutionCount is the execution count of the response.
		ExecutionCount int `json:"execution_count"`
	}
	// JupyterResult represents a response structure from Jupyter.
	JupyterResult struct {
		Text         string            `json:"text"`
		HTML         string            `json:"html"`
		Markdown     string            `json:"markdown"`
		SVG          string            `json:"svg"`
		PNG          string            `json:"png"`
		JPEG         string            `json:"jpeg"`
		PDF          string            `json:"pdf"`
		LaTeX        string            `json:"latex"`
		JSON         string            `json:"json"`
		JavaScript   string            `json:"javascript"`
		Extra        map[string]string `json:"extra"`
		IsMainResult bool              `json:"is_main_result"`
		raw          map[mime.Type]string
	}
	// Execution represents a Execution response structure from Jupyter..
	Execution struct {
		Results        []JupyterResult `json:"results"`
		Logs           Logs            `json:"logs"`
		Error          *APIError       `json:"error"`
		ExecutionCount int             `json:"execution_count"`
	}
	// APIError represents an error response structure from Jupyter.
	APIError struct {
		Name      string
		Value     string
		Traceback []string
	}
	// Logs represents a Logs response structure from Jupyter.
	Logs struct {
		// LogLevel is the log level of the logs.
		LogLevel string `json:"log_level"`
		// LogMessage is the log message of the logs.
		LogMessage string `json:"log_message"`
	}
)

// ExecuteCell executes a cell of code.
func (k *Kernel) ExecuteCell(
	ctx context.Context,
	code string,
	stdOut io.Reader,
	stdErr io.Reader,
) (ExecuteCellResponse, error) {
	return ExecuteCellResponse{}, nil
}

// Shutdown shutdowns a kernel.
func (k *Kernel) Shutdown() error {
	return nil
}

// Restart restarts a kernel.
func (k *Kernel) Restart() error {
	return nil
}
