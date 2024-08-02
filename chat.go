package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ChatRequest is a request to the chat endpoint.
type ChatRequest struct {
	Messages  []Message `json:"messages"`
	Model     string    `json:"model"`
	TopP      float64   `json:"top_p"`
	MaxTokens int       `json:"max_tokens"`
	Stop      []string  `json:"stop,omitempty"`
	Seed      int       `json:"seed"`
	Stream    bool      `json:"stream"`
	Format    struct {
		Type Format `json:"type"`
	} `json:"response_format"`
}

// Message is a message in a chat request.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse is a response from the chat endpoint.
type ChatResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int     `json:"prompt_tokens"`
		PromptTime       float64 `json:"prompt_time"`
		CompletionTokens int     `json:"completion_tokens"`
		CompletionTime   float64 `json:"completion_time"`
		ToTalTokens      int     `json:"to tal_tokens"`
		TotalTime        float64 `json:"total_time"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
	XGroq             struct {
		Id string `json:"id"`
	} `json:"x_groq"`
}

// Chat sends a request to the chat endpoint of the Groq API.
func (c *Client) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			req.Stream = false
			if !c.models.contains(req.Model) {
				return nil, fmt.Errorf("model %s is not available to the client", req.Model)
			}
			url := "https://api.groq.com/openai/v1/chat/completions"
			reqJsonBody, err := json.Marshal(req)
			if err != nil {
				return nil, fmt.Errorf("error marshalling request body: %v", err)
			}
			request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(reqJsonBody)))
			if err != nil {
				return nil, fmt.Errorf("error reading request. %v", err)
			}
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.groqApiKey))
			done, err := c.client.Do(request)
			if err != nil {
				return nil, fmt.Errorf("error sending request: %v", err)
			}
			if done.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("unexpected status code: %d", done.StatusCode)
			}
			defer done.Body.Close()
			var chatResponse ChatResponse
			err = json.NewDecoder(done.Body).Decode(&chatResponse)
			if err != nil {
				return nil, fmt.Errorf("error decoding response: %v", err)
			}
			return &chatResponse, nil
		}
	}
}

// TODO: Wait for iterators in 1.23 go
// // ChatStream sends a request to the chat stream endpoint of the Groq API.
// func (c *Client) ChatStream(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
//         req.Stream = true
//         if !c.models.contains(req.Model) {
//                 return nil, fmt.Errorf("model %s is not available to the client", req.Model)
//         }
//         url := "https://api.groq.com/openai/v1/chat/completions"
//         wr := bytes.NewBuffer([]byte{})
//         err := json.NewEncoder(wr).Encode(req)
//         if err != nil {
//                 return nil, fmt.Errorf("error marshalling request body: %v", err)
//         }
//         request, err := http.NewRequest("POST", url, wr)
//         if err != nil {
//                 return nil, fmt.Errorf("error reading request. %v", err)
//         }
//         request.Header.Set("Content-Type", "application/json")
//         request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.groqApiKey))
//         done, err := c.client.Do(request)
//         if err != nil {
//                 return nil, fmt.Errorf("error sending request: %v", err)
//         }
//         defer done.Body.Close()
//
//         if done.StatusCode != http.StatusOK {
//                 return nil, fmt.Errorf("unexpected status code: %d", done.StatusCode)
//         }
//
//         scanner := bufio.NewScanner(done.Body)
//         for scanner.Scan() {
//                 var chatResponse ChatStreamResponse
//                 err := json.Unmarshal(scanner.Bytes(), &chatResponse)
//                 if err != nil {
//                         return nil, fmt.Errorf("error decoding response: %v", err)
//                 }
//         }
//         if err := scanner.Err(); err != nil {
//                 fmt.Println("Error reading response:", err)
//         }
//         return nil, nil
// }
//
// // ChatStreamResponse is a response from the chat stream endpoint.
// type ChatStreamResponse struct {
//         ID                string `json:"id"`
//         Object            string `json:"object"`
//         Created           int    `json:"created"`
//         Model             string `json:"model"`
//         SystemFingerprint string `json:"system_fingerprint"`
//         Choices           []struct {
//                 Index int `json:"index"`
//                 Delta struct {
//                         Content string `json:"content"`
//                 } `json:"delta"`
//                 Logprobs     interface{} `json:"logprobs"`
//                 FinishReason interface{} `json:"finish_reason"`
//         } `json:"choices"`
// }
