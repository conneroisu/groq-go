package groq_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	groq "github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/internal/test"
)

func setupOpenAITestServer() (
	client *groq.Client,
	server *test.ServerTest,
	teardown func(),
) {
	server = test.NewTestServer()
	ts := server.OpenAITestServer()
	ts.Start()
	teardown = ts.Close
	config := groq.DefaultConfig(test.GetTestToken())
	config.BaseURL = ts.URL + "/v1"
	client = groq.NewClientWithConfig(config)
	return
}

// TestModeration Tests the moderations endpoint of the API using the mocked server.
func TestModerations(t *testing.T) {
	client, server, teardown := setupOpenAITestServer()
	defer teardown()
	server.RegisterHandler("/v1/moderations", handleModerationEndpoint)
	_, err := client.Moderations(context.Background(), groq.ModerationRequest{
		Model: groq.ModerationTextStable,
		Input: "I want to kill them.",
	})
	a.NoError(t, err, "Moderation error")
}

// TestModerationsWithIncorrectModel Tests passing valid and invalid models to moderations endpoint.
func TestModerationsWithDifferentModelOptions(t *testing.T) {
	var modelOptions []struct {
		model  string
		expect error
	}
	modelOptions = append(modelOptions,
		getModerationModelTestOption(groq.GPT3Dot5Turbo, groq.ErrModerationInvalidModel),
		getModerationModelTestOption(groq.ModerationTextStable, nil),
		getModerationModelTestOption(groq.ModerationTextLatest, nil),
		getModerationModelTestOption("", nil),
	)
	client, server, teardown := setupOpenAITestServer()
	defer teardown()
	server.RegisterHandler("/v1/moderations", handleModerationEndpoint)
	for _, modelTest := range modelOptions {
		_, err := client.Moderations(context.Background(), groq.ModerationRequest{
			Model: modelTest.model,
			Input: "I want to kill them.",
		})
		a.ErrorIs(t, err, modelTest.expect,
			fmt.Sprintf("Moderations(..) expects err: %v, actual err:%v", modelTest.expect, err))
	}
}

func getModerationModelTestOption(model string, expect error) struct {
	model  string
	expect error
} {
	return struct {
		model  string
		expect error
	}{model: model, expect: expect}
}

// handleModerationEndpoint Handles the moderation endpoint by the test server.
func handleModerationEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error
	var resBytes []byte

	// completions only accepts POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	var moderationReq groq.ModerationRequest
	if moderationReq, err = getModerationBody(r); err != nil {
		http.Error(w, "could not read request", http.StatusInternalServerError)
		return
	}

	resCat := groq.ResultCategories{}
	resCatScore := groq.ResultCategoryScores{}
	switch {
	case strings.Contains(moderationReq.Input, "hate"):
		resCat = groq.ResultCategories{Hate: true}
		resCatScore = groq.ResultCategoryScores{Hate: 1}

	case strings.Contains(moderationReq.Input, "hate more"):
		resCat = groq.ResultCategories{HateThreatening: true}
		resCatScore = groq.ResultCategoryScores{HateThreatening: 1}

	case strings.Contains(moderationReq.Input, "harass"):
		resCat = groq.ResultCategories{Harassment: true}
		resCatScore = groq.ResultCategoryScores{Harassment: 1}

	case strings.Contains(moderationReq.Input, "harass hard"):
		resCat = groq.ResultCategories{Harassment: true}
		resCatScore = groq.ResultCategoryScores{HarassmentThreatening: 1}

	case strings.Contains(moderationReq.Input, "suicide"):
		resCat = groq.ResultCategories{SelfHarm: true}
		resCatScore = groq.ResultCategoryScores{SelfHarm: 1}

	case strings.Contains(moderationReq.Input, "wanna suicide"):
		resCat = groq.ResultCategories{SelfHarmIntent: true}
		resCatScore = groq.ResultCategoryScores{SelfHarm: 1}

	case strings.Contains(moderationReq.Input, "drink bleach"):
		resCat = groq.ResultCategories{SelfHarmInstructions: true}
		resCatScore = groq.ResultCategoryScores{SelfHarmInstructions: 1}

	case strings.Contains(moderationReq.Input, "porn"):
		resCat = groq.ResultCategories{Sexual: true}
		resCatScore = groq.ResultCategoryScores{Sexual: 1}

	case strings.Contains(moderationReq.Input, "child porn"):
		resCat = groq.ResultCategories{SexualMinors: true}
		resCatScore = groq.ResultCategoryScores{SexualMinors: 1}

	case strings.Contains(moderationReq.Input, "kill"):
		resCat = groq.ResultCategories{Violence: true}
		resCatScore = groq.ResultCategoryScores{Violence: 1}

	case strings.Contains(moderationReq.Input, "corpse"):
		resCat = groq.ResultCategories{ViolenceGraphic: true}
		resCatScore = groq.ResultCategoryScores{ViolenceGraphic: 1}
	}

	result := groq.Result{Categories: resCat, CategoryScores: resCatScore, Flagged: true}

	res := groq.ModerationResponse{
		ID:    strconv.Itoa(int(time.Now().Unix())),
		Model: moderationReq.Model,
	}
	res.Results = append(res.Results, result)

	resBytes, _ = json.Marshal(res)
	fmt.Fprintln(w, string(resBytes))
}

// getModerationBody Returns the body of the request to do a moderation.
func getModerationBody(r *http.Request) (groq.ModerationRequest, error) {
	moderation := groq.ModerationRequest{}
	// read the request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return groq.ModerationRequest{}, err
	}
	err = json.Unmarshal(reqBody, &moderation)
	if err != nil {
		return groq.ModerationRequest{}, err
	}
	return moderation, nil
}
