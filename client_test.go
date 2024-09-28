package groq

import (
	"log/slog"
	"net/http"
	"testing"
)

func TestClient(t *testing.T) {
	client, err := NewClient(
		"test",
		WithBaseURL("http://localhost/v1"),
		WithClient(http.DefaultClient),
		WithLogger(slog.Default()),
	)
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("client is nil")
	}
}
