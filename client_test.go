package groq

import (
	"net/http"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func TestClient(t *testing.T) {
	client, err := NewClient(
		"test",
		WithBaseURL("http://localhost/v1"),
		WithClient(http.DefaultClient),
		WithLogger(zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()),
	)
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("client is nil")
	}
}
