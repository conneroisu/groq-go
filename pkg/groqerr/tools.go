package groqerr

import "fmt"

type (
	// ErrToolArgument is returned when an argument is invalid.
	ErrToolArgument struct {
		ToolName string
		ArgName  string
	}
)

// Error implements the error interface for ErrToolArgument.
func (e ErrToolArgument) Error() string {
	return fmt.Sprintf("invalid argument %s for tool %s", e.ArgName, e.ToolName)
}
