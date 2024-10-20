package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/extensions/composio"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
) error {
	key := os.Getenv("COMPOSIO_API_KEY")
	client, err := groq.NewClient(key)
	if err != nil {
		return err
	}
	comp, err := composio.NewComposer(
		key,
		composio.WithLogger(slog.Default()),
	)
	if err != nil {
		return err
	}
	tools, err := comp.GetTools(
		composio.WithApp(composio.AppGithub),
	)
	if err != nil {
		return err
	}
	chat, err := client.CreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model: groq.ModelLlama3Groq70B8192ToolUsePreview,
		Messages: []groq.ChatCompletionMessage{
			{
				Role: groq.ChatMessageRoleUser,
				Content: `
You are a github star bot.
You will be given a repo name and you will star it. 
Star a repo conneroisu/groq-go on GitHub
`,
			},
		},
		MaxTokens: 2000,
		Tools:     tools,
	})
	if err != nil {
		return err
	}
	comp.Run(ctx, chat)
	return nil
}
