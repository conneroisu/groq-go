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

func TestAuth(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	ts := test.NewTestServer()
	ts.RegisterHandler("/v1/connectedAccounts", func(w http.ResponseWriter, _ *http.Request) {
		var items struct {
			Items []composio.ConnectedAccount `json:"items"`
		}
		items.Items = append(items.Items, composio.ConnectedAccount{
			IntegrationID:      "INTEGRATION_ID",
			ID:                 "ID",
			MemberID:           "MEMBER_ID",
			ClientUniqueUserID: "CLIENT_UNIQUE_USER_ID",
			Status:             "STATUS",
			AppUniqueID:        "APP_UNIQUE_ID",
			AppName:            "APP_NAME",
			InvocationCount:    "INVOCATION_COUNT",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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
	ca, err := client.GetConnectedAccounts(ctx)
	a.NoError(err)
	assert.NotEmpty(t, ca)
}

// TestUnitGetConnectedAccounts is an Unit test using a real composio server and api key.
func TestUnitGetConnectedAccounts(t *testing.T) {
	if !test.IsIntegrationTest() {
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
