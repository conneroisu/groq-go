// Package main shows an example of using the e2b extension.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/extensions/e2b"
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
	key := os.Getenv("GROQ_KEY")
	e2bKey := os.Getenv("E2B_API_KEY")
	client, err := groq.NewClient(key)
	if err != nil {
		return err
	}
	sb, err := e2b.NewSandbox(
		ctx,
		e2bKey,
	)
	if err != nil {
		return err
	}
	defer func() {
		err := sb.Stop(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()

	chat, err := client.CreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model: groq.ModelLlama3Groq70B8192ToolUsePreview,
		Messages: []groq.ChatCompletionMessage{
			{
				Role: groq.ChatMessageRoleUser,
				Content: `

Given the tools given to you, create a golang project with the following files:

<files>
main.go
utils.go
<files>

The main function should call the "utils.run() error" function.

The project should, when run, print the following to stdout:

<output>
Hello, World!
<output>
`,
			},
		},
		MaxTokens: 2000,
		Tools:     sb.GetTools(),
	})
	if err != nil {
		return err
	}
	fmt.Println(chat.Choices[0].Message.Content)
	resp, err := sb.RunTooling(ctx, chat)
	if err != nil {
		return err
	}
	fmt.Println(resp[0].Content)
	return nil
}
