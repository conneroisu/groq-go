package composio_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/extensions/composio"
	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/conneroisu/groq-go/pkg/tools"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
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
	ts.RegisterHandler("/v2/actions/TOOL/execute", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`response1`))
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
	ca, err := client.GetConnectedAccounts(ctx, composio.WithShowActiveOnly(true))
	a.NoError(err)
	resp, err := client.Run(ctx, ca[0], groq.ChatCompletionResponse{
		Choices: []groq.ChatCompletionChoice{{
			Message: groq.ChatCompletionMessage{
				Role:    groq.RoleUser,
				Content: "Hello!",
				ToolCalls: []tools.ToolCall{{
					Function: tools.FunctionCall{
						Name:      "TOOL",
						Arguments: `{ "foo": "bar", }`,
					}}}},
			FinishReason: groq.ReasonFunctionCall,
		}}})
	a.NoError(err)
	assert.Equal(t, "response1", resp[0].Content)
}

func TestUnitRun(t *testing.T) {
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
	ts, err := client.GetTools(
		ctx, composio.WithApp("GITHUB"), composio.WithUseCase("StarRepo"))
	a.NoError(err)
	a.NotEmpty(ts)
	groqClient, err := groq.NewClient(
		os.Getenv("GROQ_KEY"),
	)
	a.NoError(err, "NewClient error")
	response, err := groqClient.ChatCompletion(ctx, groq.ChatCompletionRequest{
		Model: groq.ModelLlama3Groq8B8192ToolUsePreview,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.RoleUser,
				Content: "Star the facebookresearch/spiritlm repository on GitHub",
			},
		},
		MaxTokens: 2000,
		Tools:     ts,
	})
	a.NoError(err)
	a.NotEmpty(response.Choices[0].Message.ToolCalls)
	users, err := client.GetConnectedAccounts(ctx)
	a.NoError(err)
	resp2, err := client.Run(ctx, users[0], response)
	a.NoError(err)
	a.NotEmpty(resp2)
	t.Logf("%+v\n", resp2)
}
