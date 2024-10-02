package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/extensions/toolhouse"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// run runs the main function.
func run() error {
	ctx := context.Background()
	ext, err := toolhouse.NewExtension(os.Getenv("TOOLHOUSE_API_KEY"))
	if err != nil {
		return err
	}
	client, err := groq.NewClient(os.Getenv("GROQ_KEY"))
	if err != nil {
		return err
	}
	ts := ext.MustGetTools(ctx, toolhouse.WithBundle(
		"default",
	), toolhouse.WithMetadata(map[string]any{
		"id":       "conner",
		"timezone": 5,
	},
	), toolhouse.WithProvider(
		"openai",
	))
	ats := []groq.Tool{}
	for _, t := range ts {
		jsV, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsV))
		if t.Type == "" {
			continue
		}
		ats = append(ats, t)
	}
	re, err := client.CreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model: groq.ModelLlama3Groq70B8192ToolUsePreview,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.ChatMessageRoleUser,
				Content: "Write a python function to print the first 10 prime numbers then respond with the answer.",
			},
		},
		Tools: ats,
	})
	if err != nil {
		return fmt.Errorf("failed to create chat completion: %w", err)
	}
	fmt.Println(re.Choices[0].Message.Content)
	r, err := ext.Run(ctx, re)
	if err != nil {
		return fmt.Errorf("failed to run tool: %w", err)
	}
	fmt.Println("ran tool got:", r)
	finalr, err := client.CreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model:     groq.ModelLlama3Groq70B8192ToolUsePreview,
		Messages:  r,
		MaxTokens: 2000,
	})
	if err != nil {
		return fmt.Errorf("failed to create chat completion: %w", err)
	}
	fmt.Println(finalr.Choices[0].Message.Content)

	return nil
}
