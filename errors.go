package gogroq

// APIError provides error information returned by the Groq API.
type APIError struct {
	Code           any     `json:"code,omitempty"`
	Message        string  `json:"message"`
	Param          *string `json:"param,omitempty"`
	Type           string  `json:"type"`
	HTTPStatusCode int     `json:"-"`
}

// RequestError provides information about generic request errors.
type RequestError struct {
	HTTPStatusCode int
	Err            error
}

// ErrorResponse provides information about errors returned by the Groq API.
type ErrorResponse struct {
	Error *APIError `json:"error,omitempty"`
}

// Error returns the error message of an APIError.
func (e *APIError) Error() string {
	return e.Message
}

// Error returns the error message of a RequestError.
func (e *RequestError) Error() string {
	return e.Err.Error()
}
