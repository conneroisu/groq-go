package composio

import (
	"context"
	"testing"

	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
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
	ts, err := client.GetConnectedAccounts(ctx)
	a.NoError(err)
	a.NotEmpty(ts)
}
