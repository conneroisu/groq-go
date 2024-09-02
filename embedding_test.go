// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package groq_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/internal/testutil"
	"github.com/conneroisu/groq-go/option"
	"github.com/conneroisu/groq-go/shared"
)

func TestEmbeddingNewWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := groq.NewClient(
		option.WithBaseURL(baseURL),
		option.WithBearerToken("My Bearer Token"),
	)
	_, err := client.Embeddings.New(context.TODO(), groq.EmbeddingNewParams{
		Input:          groq.F[groq.EmbeddingNewParamsInputUnion](shared.UnionString("The quick brown fox jumped over the lazy dog")),
		Model:          groq.F(groq.EmbeddingNewParamsModelNomicEmbedTextV1_5),
		EncodingFormat: groq.F(groq.EmbeddingNewParamsEncodingFormatFloat),
		User:           groq.F("user"),
	})
	if err != nil {
		var apierr *groq.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
