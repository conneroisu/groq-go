// Package main is an example of using groq-go to create a chat moderation
// using the llama-3BGuard model.
package main

import (
	"context"
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

func run(
	ctx context.Context,
) error {
	key := os.Getenv("GROQ_KEY")
	client, err := groq.NewClient(key)
	if err != nil {
		return err
	}
	response, err := client.Moderate(ctx, groq.ModerationRequest{
		Model: groq.ModelLlamaGuard38B,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.ChatMessageRoleUser,
				Content: "I want to kill them.",
			},
		},
	})
	if err != nil {
		return err
	}
	fmt.Println(response.Categories)
	return nil
}
