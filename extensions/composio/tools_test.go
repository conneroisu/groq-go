package composio

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

// TestGetTools tests the ability of the composio client to get tools.
func TestGetTools(t *testing.T) {
	if !test.IsUnitTest() {
		t.Skip()
	}
	a := assert.New(t)
	ctx := context.Background()
	key, err := test.GetAPIKey("COMPOSIO_API_KEY")
	a.NoError(err)
	client, err := NewComposer(
		key,
		WithLogger(test.DefaultLogger),
	)
	a.NoError(err)
	ts, err := client.GetTools(ctx)
	a.NoError(err)
	a.NotEmpty(ts)
	jsval, err := json.MarshalIndent(ts, "", "  ")
	a.NoError(err)
	f, err := os.Create("tools.json")
	a.NoError(err)
	defer f.Close()
	_, err = f.Write(jsval)
	a.NoError(err)
}
