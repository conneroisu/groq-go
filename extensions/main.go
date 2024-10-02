package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/extensions/toolhouse"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// run runs the main function.
func run() error {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	ctx := context.Background()
	ext, err := toolhouse.NewExtension(os.Getenv("TOOLHOUSE_API_KEY"))
	if err != nil {
		return err
	}
	client, err := groq.NewClient(os.Getenv("GROQ_KEY"))
	if err != nil {
		return err
	}
	history := []groq.ChatCompletionMessage{
		{
			Role:    groq.ChatMessageRoleUser,
			Content: "Write a python function to print the first 10 prime numbers containing the number 3 then respond with the answer. DO NOT GUESS WHAT THE OUTPUT SHOULD BE. MAKE SURE TO CALL THE TOOL GIVEN.",
		},
	}
	re, err := client.CreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model:    groq.ModelLlama3Groq70B8192ToolUsePreview,
		Messages: history,
		Tools: ext.MustGetTools(ctx, toolhouse.WithMetadata(map[string]any{
			"id":       "conner",
			"timezone": 5,
		})),
		ToolChoice: "required",
	})
	if err != nil {
		return fmt.Errorf("failed to create 1 chat completion: %w", err)
	}
	history = append(history, re.Choices[0].Message)
	r, err := ext.Run(ctx, re, map[string]any{
		"id":       "conner",
		"timezone": 5,
	})
	if err != nil {
		return fmt.Errorf("failed to run tool: %w", err)
	}
	history = append(history, r...)
	finalr, err := client.CreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model:     groq.ModelLlama3Groq70B8192ToolUsePreview,
		Messages:  history,
		MaxTokens: 2000,
	})
	if err != nil {
		return fmt.Errorf("failed to create 2 chat completion: %w", err)
	}
	history = append(history, finalr.Choices[0].Message)
	jsnHistory, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}
	fmt.Println(string(jsnHistory))
	return nil
}

func printHistoryContents(history []groq.ChatCompletionMessage) {
	for _, m := range history {
		fmt.Println(m.Content)
	}
}
