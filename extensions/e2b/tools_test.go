package e2b

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/conneroisu/groq-go"
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

func getapiKey(t *testing.T, val string) string {
	apiKey := os.Getenv(val)
	if apiKey == "" {
		t.Fail()
	}
	return apiKey
}

func TestSandboxTooling(t *testing.T) {
	if os.Getenv("UNIT") == "" {
		t.Skip("Skipping Tooling test")
	}
	a := assert.New(t)
	ctx := context.Background()
	sb, err := NewSandbox(
		ctx,
		getapiKey(t, "E2B_API_KEY"),
		WithLogger(defaultLogger),
		WithCwd("/code"),
	)
	a.NoError(err, "NewSandbox error")
	client, err := groq.NewClient(getapiKey(t, "GROQ_KEY"))
	a.NoError(err, "NewClient error")

	tts, err := sb.getTools()
	a.NoError(err)
	// ask the ai to create a file with the data "Hello World!" in file "hello.txt"
	response, err := client.CreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model: groq.ModelLlama3Groq70B8192ToolUsePreview,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.ChatMessageRoleUser,
				Content: "Create a file called hello.txt with the data Hello World! NOTE: You are in the correct directory.",
			},
		},
		MaxTokens: 2000,
		Tools:     tts,
	})
	a.NoError(err)
	sb.logger.Debug("response from model", "response", response)
	resps, err := sb.RunTooling(ctx, response)
	a.NoError(err)
	sb.logger.Debug("tooling response", "response", resps)
	lsres, err := sb.Ls(".")
	a.NoError(err)
	a.Contains(lsres, LsResult{
		Name:  "hello.txt",
		IsDir: false,
	})
}
