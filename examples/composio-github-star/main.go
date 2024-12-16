// Package main is an example of using the composio client.
//
// It shows how to use the composio client to star a github repository.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/extensions/composio"
	"github.com/conneroisu/groq-go/pkg/test"
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
	key, err := test.GetAPIKey("GROQ_KEY")
	if err != nil {
		return err
	}
	client, err := groq.NewClient(key)
	if err != nil {
		return err
	}
	key, err = test.GetAPIKey("COMPOSIO_API_KEY")
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
		ctx,
		composio.WithApp("GITHUB"),
		composio.WithUseCase("star-repo"),
	)
	if err != nil {
		return err
	}
	chat, err := client.ChatCompletion(ctx, groq.ChatCompletionRequest{
		Model: groq.ModelLlama3Groq70B8192ToolUsePreview,
		Messages: []groq.ChatCompletionMessage{
			{
				Role: groq.RoleUser,
				Content: `
You are a github star bot. You will be given a repo name and you will star it. 
Star the repo conneroisu/groq-go on GitHub.
`,
			},
		},
		MaxTokens: 2000,
		Tools:     tools,
	})
	if err != nil {
		return err
	}
	user, err := comp.GetConnectedAccounts(ctx)
	if err != nil {
		return err
	}
	resp, err := comp.Run(ctx, user[0], chat)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}
