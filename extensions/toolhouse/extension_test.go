package toolhouse_test

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/extensions/toolhouse"
	"github.com/stretchr/testify/assert"
)

var (
	defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{}
			}
			if a.Key == "level" {
				return slog.Attr{}
			}
			if a.Key == "source" {
				str := a.Value.String()
				split := strings.Split(str, "/")
				if len(split) > 2 {
					a.Value = slog.StringValue(strings.Join(split[len(split)-2:], "/"))
				}
			}
			return a
		}}))
)

func TestNewExtension(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	if os.Getenv("UNIT") == "" {
		t.Skip("Skipping Toolhouse extension test")
	}

	ext, err := toolhouse.NewExtension(os.Getenv("TOOLHOUSE_API_KEY"),
		toolhouse.WithMetadata(map[string]any{
			"id":       "conner",
			"timezone": 5,
		}),
		toolhouse.WithLogger(defaultLogger),
	)
	a.NoError(err)
	client, err := groq.NewClient(os.Getenv("GROQ_KEY"))
	a.NoError(err)
	history := []groq.ChatCompletionMessage{
		{
			Role:    groq.ChatMessageRoleUser,
			Content: "Write a python function to print the first 10 prime numbers containing the number 3 then respond with the answer. DO NOT GUESS WHAT THE OUTPUT SHOULD BE. MAKE SURE TO CALL THE TOOL GIVEN.",
		},
	}
	print(history[len(history)-1].Content)
	re, err := client.CreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model:      groq.ModelLlama3Groq70B8192ToolUsePreview,
		Messages:   history,
		Tools:      ext.MustGetTools(ctx),
		ToolChoice: "required",
	})
	a.NoError(err)
	history = append(history, re.Choices[0].Message)
	print(history[len(history)-1].ToolCalls[len(history[len(history)-1].ToolCalls)-1].Function.Arguments)
	r, err := ext.Run(ctx, re)
	a.NoError(err)
	history = append(history, r...)
	finalr, err := client.CreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model:     groq.ModelLlama3Groq70B8192ToolUsePreview,
		Messages:  history,
		MaxTokens: 2000,
	})
	a.NoError(err)
	history = append(history, finalr.Choices[0].Message)
	print(history[len(history)-1].Content)
}
