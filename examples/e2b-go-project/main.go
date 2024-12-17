// Package main shows an example of using the e2b extension.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/extensions/e2b"
	"github.com/conneroisu/groq-go/pkg/tools"
)

var (
	history = []groq.ChatCompletionMessage{
		{
			Role: groq.RoleUser,
			Content: `
Given the callable tools provided, create a python project with the following files:

<files>
main.py
utils.py
<files>

The main function should call the "utils.run()" function.

The project should, when run, print the following to stdout:

<output>
Hello, World!
<output>

You should finish with the following shell command:

<shell-command>
python main.py
</shell-command>
`,
		},
	}
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
	groqKey := os.Getenv("GROQ_KEY")
	e2bKey := os.Getenv("E2B_API_KEY")
	client, err := groq.NewClient(groqKey)
	if err != nil {
		return err
	}
	sb, err := e2b.NewSandbox(ctx, e2bKey)
	if err != nil {
		return err
	}
	defer func() {
		err := sb.Stop(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
	ts := sb.GetTools()
	ts = append(ts, tools.Tool{
		Type: tools.ToolTypeFunction,
		Function: tools.FunctionDefinition{
			Name:        "complete",
			Description: "Signify that the assigned task is complete.",
			Parameters: tools.FunctionParameters{
				Type: "object",
				Properties: map[string]tools.PropertyDefinition{
					"task": {
						Type:        "string",
						Description: "The task that is complete.",
					}}}}})
	for {
		chat, err := client.ChatCompletion(ctx, groq.ChatCompletionRequest{
			Model:     groq.ModelLlama3Groq8B8192ToolUsePreview,
			Messages:  history,
			MaxTokens: 3000,
			Tools:     ts,
		})
		if err != nil {
			return err
		}
		if chat.Choices[0].FinishReason == groq.ReasonFunctionCall {
			if chat.Choices[0].Message.FunctionCall.Name == "complete" {
				break
			}
		}
		resp, err := sb.RunTooling(ctx, chat)
		if err != nil {
			history = append(history,
				groq.ChatCompletionMessage{
					Role:    groq.RoleUser,
					Content: err.Error(),
				})
			continue
		}
		history = append(history, resp...)
	}
	return nil
}
