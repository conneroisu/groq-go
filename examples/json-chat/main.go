// Package main demonstrates an example application of groq-go.
// It shows how to use groq-go to create a chat completion of a json object
// using the llama-3.1-70B-8192-tool-use-preview model.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/conneroisu/groq-go"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Responses is a response from the models endpoint.
type Responses []struct {
	Title string `json:"title" jsonschema:"title=Poem Title,description=Title of the poem, minLength=1, maxLength=20"`
	Text  string `json:"text" jsonschema:"title=Poem Text,description=Text of the poem, minLength=10, maxLength=200"`
}

func run(
	ctx context.Context,
) error {
	client, err := groq.NewClient(os.Getenv("GROQ_KEY"))
	if err != nil {
		return err
	}
	resp := &Responses{}
	err = client.CreateChatCompletionJSON(ctx, groq.ChatCompletionRequest{
		Model: groq.ModelLlama3Groq70B8192ToolUsePreview,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.RoleUser,
				Content: "Create 5 short poems in json format with title and text.",
			},
		},
		MaxTokens: 2000,
	}, resp)
	if err != nil {
		return err
	}

	jsValue, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsValue))

	return nil
}
