package jigsawstack_test

import (
	"context"
	"testing"

	"github.com/conneroisu/groq-go/extensions/jigsawstack"
	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

// TestJigsawStack_TextToSQL tests the TextToSQL method of the JigsawStack client.
func TestJigsawStack_TextToSQL(t *testing.T) {
	if !test.IsIntegrationTest() {
		t.Skip("Skipping unit test")
	}
	a := assert.New(t)
	ctx := context.Background()
	apiKey, err := test.GetAPIKey("JIGSAWSTACK_API_KEY")
	a.NoError(err)
	j, err := jigsawstack.NewJigsawStack(apiKey)
	a.NoError(err)
	resp, err := j.TextToSQL(ctx, "select all users", `
CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(255),
  email VARCHAR(255),
  age INT
);
`)
	a.NoError(err)
	a.NotEmpty(resp.SQL)
}
