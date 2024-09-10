// Package main demonstrates an example application of groq-go.
// It shows how to use groq-go to create a chat completion of a json object
// using the llama-3.1-8b-instant model.
package main

// url: https://cdnimg.webstaurantstore.com/images/products/large/87539/251494.jpg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/conneroisu/groq-go"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Response is a response from the models endpoint.
type Response struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
type resps []Response

func (r *resps) MarshalJSON() ([]byte, error) {
	return json.Marshal(r)
}
func run(
	ctx context.Context,
) error {
	key := os.Getenv("GROQ_KEY")
	client, err := groq.NewClient(key)
	if err != nil {
		return err
	}
	resp := resps{}
	err = client.CreateChatCompletionJSON(ctx, groq.ChatCompletionRequest{
		Model: groq.LlavaV157B4096Preview,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.ChatMessageRoleUser,
				Content: "Create 5 short poems in json format with title and text.",
			},
		},
		MaxTokens: 2000,
	}, resp)
	if err != nil {
		return err
	}

	fmt.Println(resp)

	return nil
}
