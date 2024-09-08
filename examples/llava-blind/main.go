// Package main demonstrates an example application of groq-go.
package main

// url: https://cdnimg.webstaurantstore.com/images/products/large/87539/251494.jpg

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

	response, err := client.CreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model: groq.LlavaV157B4096Preview,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.ChatMessageRoleSystem,
				Content: "You are a helpful assistant asked to identify the contents of user given image..",
			},
			{
				Role: groq.ChatMessageRoleUser,
				MultiContent: []groq.ChatMessagePart{
					{
						Type:     groq.ChatMessagePartTypeImageURL,
						ImageURL: &groq.ChatMessageImageURL{URL: "https://cdnimg.webstaurantstore.com/images/products/large/87539/251494.jpg"},
					},
				},
			},
		},
		MaxTokens: 2000,
	})
	if err != nil {
		return err
	}

	fmt.Println(response.Choices[0].Message.Content)

	return nil
}
