package test

import (
	"log"
	"testing"

	groq "github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/internal/test"
	"github.com/stretchr/testify/assert"
)

func setupGroqTestServer() (
	client *groq.Client,
	server *test.ServerTest,
	teardown func(),
) {
	server = test.NewTestServer()
	ts := server.GroqTestServer()
	ts.Start()
	teardown = ts.Close
	client, err := groq.NewClient(
		test.GetTestToken(),
		groq.WithBaseURL(ts.URL+"/v1"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func TestEmptyKeyClientCreation(t *testing.T) {
	client, err := groq.NewClient("")
	a := assert.New(t)
	a.Error(err, "NewClient should return error")
	a.Nil(client, "NewClient should return nil")
}
