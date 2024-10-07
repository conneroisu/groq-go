package e2b

import (
	"context"
	"io"
)

// ExecuteCell executes a cell of code.
func (e *Extension) ExecuteCell(
	ctx context.Context,
	code string,
	stdOut io.Reader,
	stdErr io.Reader,
) (ExecuteCellResponse, error) {
	return ExecuteCellResponse{}, nil
}

type (
	// ExecuteCellResponse represents the response of the execute cell method.
	ExecuteCellResponse struct {
		// CellID is the cell id of the response.
		CellID string `json:"cell_id"`
		// ExecutionCount is the execution count of the response.
		ExecutionCount int `json:"execution_count"`
	}
)
