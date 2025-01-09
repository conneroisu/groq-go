package jigsawstack_test

import (
	"context"
	"testing"

	"github.com/conneroisu/groq-go/extensions/jigsawstack"
	"github.com/conneroisu/groq-go/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestJigsawStack_PromptCreate(t *testing.T) {
	if !test.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}
	a := assert.New(t)
	apiKey, err := test.GetAPIKey("JIGSAWSTACK_API_KEY")
	a.NoError(err)
	j, err := jigsawstack.NewJigsawStack(apiKey)
	a.NoError(err)
	resp, err := j.PromptCreate(context.Background(), jigsawstack.PromptCreateRequest{
		Prompt: `
You are a helpful assistant that answers questions based on the provided context.
Your job is to provide code completions based on the provided context.
		`,
		Inputs: []jigsawstack.PromptCreateInput{
			{
				Key:      "context",
				Optional: false,
				InitialValue: `
<file name="context.py">
def main():
    print("Hello, World!")

if __name__ == "__main__":
    main()
</file>
				`,
			},
		},
	})
	a.NoError(err)
	t.Logf("response: %v", resp)
	t.Fail()
}
func TestJigsawStack_PromptGet(t *testing.T) {
	if !test.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}
	a := assert.New(t)
	apiKey, err := test.GetAPIKey("JIGSAWSTACK_API_KEY")
	a.NoError(err)
	j, err := jigsawstack.NewJigsawStack(apiKey)
	a.NoError(err)
	resp, err := j.PromptGet(context.Background(), "test")
	a.NoError(err)
	a.NotEmpty(resp.Prompt)
}
