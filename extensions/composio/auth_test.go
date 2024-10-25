package composio_test

import (
	"context"
	"testing"

	"github.com/conneroisu/groq-go/extensions/composio"
	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

// TestUnitGetConnectedAccounts is an Unit test using a real composio server and api key.
func TestUnitGetConnectedAccounts(t *testing.T) {
	if !test.IsUnitTest() {
		t.Skip()
	}
	a := assert.New(t)
	ctx := context.Background()
	key, err := test.GetAPIKey("COMPOSIO_API_KEY")
	a.NoError(err)
	client, err := composio.NewComposer(
		key,
		composio.WithLogger(test.DefaultLogger),
	)
	a.NoError(err)
	ts, err := client.GetConnectedAccounts(ctx)
	a.NoError(err)
	a.NotEmpty(ts)
}
