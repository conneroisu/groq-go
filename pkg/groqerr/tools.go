package groqerr

import "fmt"

type (
	// ErrToolArgument is returned when an argument is invalid.
	ErrToolArgument struct {
		ToolName string
		ArgName  string
	}
	// ErrNonFunctionCall is returned when a expected toolcall is not a function call.
	ErrNonFunctionCall struct {
	}
)

// Error implements the error interface for ErrToolArgument.
func (e ErrToolArgument) Error() string {
	return fmt.Sprintf("invalid argument %s for tool %s", e.ArgName, e.ToolName)
}

// Error implements the error interface for ErrNonFunctionCall.
func (e ErrNonFunctionCall) Error() string {
	return "ran on response without a function call"
}
