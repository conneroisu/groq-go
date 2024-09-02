// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package groq_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/internal/testutil"
	"github.com/conneroisu/groq-go/option"
)

func TestAudioTranslationNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Audio.Translations.New(context.TODO(), groq.AudioTranslationNewParams{
		File:           groq.F(io.Reader(bytes.NewBuffer([]byte("some file contents")))),
		Model:          groq.F(groq.AudioTranslationNewParamsModelWhisperLargeV3),
		Prompt:         groq.F("prompt"),
		ResponseFormat: groq.F(groq.AudioTranslationNewParamsResponseFormatJson),
		Temperature:    groq.F(0.000000),
	})
	if err != nil {
		var apierr *groq.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
