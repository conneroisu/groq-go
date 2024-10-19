package e2b_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/conneroisu/groq-go/extensions/e2b"
	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func TestNewSandbox(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	srv := test.NewTestServer()
	ts := srv.E2bTestServer()
	wsts := httptest.NewServer(http.HandlerFunc(echo))
	srv.RegisterHandler("/sandboxes", func(w http.ResponseWriter, r *http.Request) {
		srv.Logger.Debug("received sandboxes request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"sandboxID": "test-sandbox"}`))
		if err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})
	ts.Start()
	u := "ws" + strings.TrimPrefix(wsts.URL, "http")
	sb, err := e2b.NewSandbox(
		ctx,
		test.GetTestToken(),
		e2b.WithLogger(defaultLogger),
		e2b.WithBaseURL(ts.URL),
		e2b.WithWsURL(func(s *e2b.Sandbox) string {
			return u + "/ws"
		}),
	)
	a.NoError(err, "NewSandbox error")
	a.NotNil(sb, "NewSandbox returned nil")
}
