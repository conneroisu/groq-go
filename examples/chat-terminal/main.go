// Package main demonstrates how to use groq-go to create a chat application
// using the groq api accessable through the terminal.
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/conneroisu/groq-go"
)

var (
	history = []groq.ChatCompletionMessage{}
)

func main() {
	if err := run(
		context.Background(),
		os.Stdin,
		os.Stdout,
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	r io.Reader,
	w io.Writer,
) error {
	key := os.Getenv("GROQ_KEY")
	client, err := groq.NewClient(key)
	if err != nil {
		return err
	}
	for {
		err = input(ctx, r, w, client)
		if err != nil {
			return err
		}
	}
}

func input(ctx context.Context, r io.Reader, w io.Writer, client *groq.Client) error {
	fmt.Println("")
	fmt.Print("->")
	reader := bufio.NewReader(r)
	writer := w
	var lines []string
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		if len(strings.TrimSpace(line)) == 0 {
			break
		}
		lines = append(lines, line)
		break
	}
	in := strings.Join(lines, "\n")
	history = append(history, groq.ChatCompletionMessage{
		Role:    groq.ChatMessageRoleUser,
		Content: in,
	})
	output, err := client.CreateChatCompletionStream(
		ctx,
		groq.ChatCompletionRequest{
			Model:     groq.ModelGemma29BIt,
			Messages:  history,
			MaxTokens: 2000,
		},
	)
	if err != nil {
		return err
	}
	fmt.Fprintln(writer, "")
	fmt.Fprint(writer, "ai: ")
	for {
		response, err := output.Recv()
		if err != nil {
			return err
		}
		if response.Choices[0].FinishReason == groq.FinishReasonStop {
			break
		}
		fmt.Fprint(writer, response.Choices[0].Delta.Content)
	}
	return nil
}
