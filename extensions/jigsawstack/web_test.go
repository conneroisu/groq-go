package jigsawstack_test

import (
	"context"
	"testing"

	"github.com/conneroisu/groq-go/extensions/jigsawstack"
	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

// TestJigsawStack_WebSearch tests the WebSearch method of the JigsawStack client.
func TestJigsawStack_WebSearch(t *testing.T) {
	if !test.IsIntegrationTest() {
		t.Skip("Skipping unit test")
	}
	a := assert.New(t)
	ctx := context.Background()
	apiKey, err := test.GetAPIKey("JIGSAWSTACK_API_KEY")
	a.NoError(err)
	j, err := jigsawstack.NewJigsawStack(apiKey)
	a.NoError(err)
	resp, err := j.WebSearch(ctx, "hello world golang")
	a.NoError(err)
	a.NotEmpty(resp.Results)
}

// TestJigsawStack_WebSearchSuggestions tests the WebSearchSuggestions method of the JigsawStack client.
func TestJigsawStack_WebSearchSuggestions(t *testing.T) {
	if !test.IsIntegrationTest() {
		t.Skip("Skipping unit test")
	}
	a := assert.New(t)
	ctx := context.Background()
	apiKey, err := test.GetAPIKey("JIGSAWSTACK_API_KEY")
	a.NoError(err)
	j, err := jigsawstack.NewJigsawStack(apiKey)
	a.NoError(err)
	resp, err := j.WebSearchSuggestions(ctx, "hello")
	a.NoError(err)
	a.NotEmpty(resp.Suggestions)
}
