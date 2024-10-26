package e2b

import (
	"context"
	"os"
	"testing"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

func getapiKey(t *testing.T, val string) string {
	apiKey := os.Getenv(val)
	if apiKey == "" {
		t.Fail()
	}
	return apiKey
}

func TestSandboxTooling(t *testing.T) {
	if !test.IsUnitTest() {
		t.Skip()
	}
	a := assert.New(t)
	ctx := context.Background()
	sb, err := NewSandbox(
		ctx,
		getapiKey(t, "E2B_API_KEY"),
		WithLogger(test.DefaultLogger),
		WithCwd("/code"),
	)
	a.NoError(err, "NewSandbox error")
	client, err := groq.NewClient(getapiKey(t, "GROQ_KEY"))
	a.NoError(err, "NewClient error")
	tools := sb.GetTools()
	// ask the ai to create a file with the data "Hello World!" in file "hello.txt"
	response := client.MustCreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model: groq.ModelLlama3Groq70B8192ToolUsePreview,
		Messages: []groq.ChatCompletionMessage{
			{
				Role: groq.ChatMessageRoleUser,
				Content: `
Create a file called 'hello.txt' with the data:
<file name="hello.txt">
Hello World! 
</file>
NOTE: You are in the correct cwd. Just call the write tool with a name of hello.txt and data of Hello World!
`,
			},
		},
		MaxTokens: 2000,
		Tools:     tools,
	})
	sb.logger.Debug("response from model", "response", response)
	resps, err := sb.RunTooling(ctx, response)
	a.NoError(err)
	sb.logger.Debug("tooling response", "response", resps)
	lsres, err := sb.Ls(ctx, ".")
	a.NoError(err)
	a.Contains(lsres, LsResult{
		Name:  "hello.txt",
		IsDir: false,
	})
}
