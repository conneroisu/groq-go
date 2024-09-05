package groq_test

import (
	"log"

	groq "github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/internal/test"
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
