package groq_test

import (
	"log/slog"
	"net/http"
	"testing"

	groq "github.com/conneroisu/groq-go"
	"github.com/stretchr/testify/assert"
)

// TestClient tests the creation of a new client.
func TestClient(t *testing.T) {
	a := assert.New(t)
	client, err := groq.NewClient(
		"test",
		groq.WithBaseURL("http://localhost/v1"),
		groq.WithClient(http.DefaultClient),
		groq.WithLogger(slog.Default()),
	)
	a.NoError(err)
	a.NotNil(client)
}
