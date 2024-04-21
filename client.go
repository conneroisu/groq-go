package gogroq

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("vim-go")
}

// Client is a Groq client
type Client struct {
	groqApiKey string       // Groq API key
	client     *http.Client // Client is the HTTP client to use
	BaseURL    string       // BaseURL is the base URL for the Groq API
}

// GroqOpts is a function that sets options for a Groq client
type GroqOpts func(*Client)

// WithBaseURL sets the base URL for the Groq client
func WithBaseURL(url string) GroqOpts {
	return func(c *Client) {
		c.BaseURL = url
	}
}

// NewClient creates a new Groq client
func NewClient(groqApiKey string, client *http.Client, opts ...GroqOpts) *Client {
	c := &Client{
		groqApiKey: groqApiKey,
		client:     client,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}
