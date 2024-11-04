package groqerr

import "fmt"

type (
	// ErrContentFieldsMisused is an error that occurs when both Content and
	// MultiContent properties are set.
	ErrContentFieldsMisused struct{}
	// ErrToolNotFound is returned when a tool is not found.
	ErrToolNotFound struct {
		ToolName string
	}
)

// Error implements the error interface.
func (e ErrContentFieldsMisused) Error() string {
	return fmt.Errorf("can't use both Content and MultiContent properties simultaneously").
		Error()
}

type (
	// ErrRequest is a request error.
	ErrRequest struct {
		HTTPStatusCode int
		Err            error
	}
)

// Error implements the error interface.
func (e *ErrRequest) Error() string {
	return fmt.Sprintf(
		"error, status code: %d, message: %s",
		e.HTTPStatusCode,
		e.Err,
	)
}

// Unwrap unwraps the error.
func (e *ErrRequest) Unwrap() error {
	return e.Err
}

// Error implements the error interface.
func (e ErrToolNotFound) Error() string {
	return fmt.Sprintf("tool %s not found", e.ToolName)
}
