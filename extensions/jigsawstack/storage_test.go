package jigsawstack_test

import (
	"context"
	"testing"

	"github.com/conneroisu/groq-go/extensions/jigsawstack"
	"github.com/conneroisu/groq-go/internal/test"
	"github.com/stretchr/testify/assert"
)

// TestJigsawStack_FileAdd tests the FileAdd method of the JigsawStack client.
func TestJigsawStack_FileAdd(t *testing.T) {
	if !test.IsIntegrationTest() {
		t.Skip("Skipping unit test")
	}
	a := assert.New(t)
	ctx := context.Background()
	apiKey, err := test.GetAPIKey("JIGSAWSTACK_API_KEY")
	a.NoError(err)
	j, err := jigsawstack.NewJigsawStack(apiKey)
	a.NoError(err)
	resp, err := j.FileAdd(ctx, "test", "text/plain", "hello world")
	a.NoError(err)
	a.NotEmpty(resp)
}
