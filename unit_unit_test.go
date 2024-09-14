//go:build !test
// +build !test

package groq_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"

	"github.com/conneroisu/groq-go"
	"github.com/stretchr/testify/assert"
)

func TestTestServer(t *testing.T) {
	num := rand.Intn(100)
	a := assert.New(t)
	ctx := context.Background()
	client, err := groq.NewClient(os.Getenv("GROQ_KEY"))
	a.NoError(err, "NewClient error")
	strm, err := client.CreateChatCompletionStream(
		ctx,
		groq.ChatCompletionRequest{
			Model: groq.Llama38B8192,
			Messages: []groq.ChatCompletionMessage{
				{
					Role: groq.ChatMessageRoleUser,
					Content: fmt.Sprintf(`
problem: %d
You have a six-sided die that you roll once. Let $R{i}$ denote the event that the roll is $i$. Let $G{j}$ denote the event that the roll is greater than $j$. Let $E$ denote the event that the roll of the die is even-numbered.
(a) What is $P\left[R{3} \mid G{1}\right]$, the conditional probability that 3 is rolled given that the roll is greater than 1 ?
(b) What is the conditional probability that 6 is rolled given that the roll is greater than 3 ?
(c) What is the $P\left[G_{3} \mid E\right]$, the conditional probability that the roll is greater than 3 given that the roll is even?
(d) Given that the roll is greater than 3, what is the conditional probability that the roll is even?
					`, num,
					),
				},
			},
			MaxTokens: 2000,
			Stream:    true,
		},
	)
	a.NoError(err, "CreateCompletionStream error")

	i := 0
	for {
		i++
		val, err := strm.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		// t.Logf("%d %s\n", i, val.Choices[0].Delta.Content)
		print(val.Choices[0].Delta.Content)
	}
}
