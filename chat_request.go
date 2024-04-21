package gogroq

// curl -X POST "https://api.groq.com/openai/v1/chat/completions" \
//      -H "Authorization: Bearer $GROQ_API_KEY" \
//      -H "Content-Type: application/json" \
//      -d '{"messages": [{"role": "user", "content": "Explain the importance of fast language models"}], "model": "mixtral-8x7b-32768"}'
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) newChatReq(req ChatRequest) (*http.Request, error) {
	url := fmt.Sprintf("%s/openai/v1/chat/completions", c.BaseURL)
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
	return request, nil
}

func (c *Client) parseChatResp(resp *http.Response) (*ChatResponse, error) {
	var chatResponse ChatResponse
	err := json.NewDecoder(resp.Body).Decode(&chatResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}
	return &chatResponse, nil
}
