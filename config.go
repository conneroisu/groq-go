package groq

import (
	"net/http"
)

const (
	openaiAPIURLv1                 = "https://api.openai.com/v1"
	defaultEmptyMessagesLimit uint = 300
)

// APIType is the type of API.
type APIType string

const (
	// APITypeOpenAI is the OpenAI API type.
	APITypeOpenAI APIType = "OPEN_AI"
)

// HTTPDoer is an interface for making HTTP requests.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	authToken string

	BaseURL    string
	OrgID      string
	APIType    APIType
	HTTPClient HTTPDoer

	EmptyMessagesLimit uint
}

// DefaultConfig returns a ClientConfig with default values.
func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		authToken: authToken,
		BaseURL:   openaiAPIURLv1,
		APIType:   APITypeOpenAI,
		OrgID:     "",

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}

// String returns a string representation of the ClientConfig.
func (ClientConfig) String() string {
	return "<OpenAI API ClientConfig>"
}
