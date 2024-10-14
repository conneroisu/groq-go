package e2b

import "fmt"

type (
	// ErrToolNotFound is returned when a tool is not found.
	ErrToolNotFound struct {
		ToolName string
	}
	// ErrToolArgument is returned when an argument is invalid.
	ErrToolArgument struct {
		ToolName string
		ArgName  string
	}
	// ErrMissingRequiredArgument is returned when a required argument is missing.
	ErrMissingRequiredArgument struct {
		ToolName string
		ArgName  string
	}
	ErrArgumentsUnmarshallable struct {
	}
)

// Error implements the error interface for ErrToolNotFound.
func (e ErrToolNotFound) Error() string {
	return fmt.Sprintf("tool %s not found", e.ToolName)
}

// Error implements the error interface for ErrToolArgument.
func (e ErrToolArgument) Error() string {
	return fmt.Sprintf("invalid argument %s for tool %s", e.ArgName, e.ToolName)
}

// Error implements the error interface	for ErrMissingRequiredArgument.
func (e ErrMissingRequiredArgument) Error() string {
	return fmt.Sprintf("missing required argument %s for tool %s", e.ArgName, e.ToolName)
}

func (e ErrArgumentsUnmarshallable) Error() string {
	return "given json arguments cannot be unmarshalled"
}
