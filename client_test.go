package groq

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClient tests the creation of a new client.
func TestClient(t *testing.T) {
	a := assert.New(t)
	client, err := NewClient(
		"test",
		WithBaseURL("http://localhost/v1"),
		WithClient(http.DefaultClient),
		WithLogger(slog.Default()),
	)
	a.NoError(err)
	a.NotNil(client)
}
