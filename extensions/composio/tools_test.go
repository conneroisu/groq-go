package composio_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/conneroisu/groq-go/extensions/composio"
	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestGetTools(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	ts := test.NewTestServer()
	ts.RegisterHandler("/v1/actions", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		var items struct {
			Items []composio.Tool `json:"items"`
		}
		items.Items = append(items.Items, composio.Tool{
			Enum:        "enum",
			Name:        "NAME",
			Tags:        []string{"TAG"},
			DisplayName: "DISPLAY_NAME",
		})
		jsonBytes, err := json.Marshal(items)
		a.NoError(err)
		_, err = w.Write(jsonBytes)
		a.NoError(err)
	})

	testS := ts.ComposioTestServer()
	testS.Start()
	client, err := composio.NewComposer(
		test.GetTestToken(),
		composio.WithLogger(test.DefaultLogger),
		composio.WithBaseURL(testS.URL),
	)
	a.NoError(err)
	ca, err := client.GetTools(ctx)
	a.NoError(err)
	a.NotEmpty(ca)
}

// TestUnitGetTools tests the ability of the composio client to get tools.
func TestUnitGetTools(t *testing.T) {
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
	ts, err := client.GetTools(ctx, composio.WithApp("GITHUB"))
	a.NoError(err)
	a.NotEmpty(ts)
}
