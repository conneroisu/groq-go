// Package main is an example of using the groq-go library to create a
// transcription/translation using the whisper-large-v3 model.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/models"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
) error {
	client, err := groq.NewClient(os.Getenv("GROQ_KEY"))
	if err != nil {
		return err
	}
	response, err := client.CreateTranslation(ctx, groq.AudioRequest{
		Model:    models.ModelWhisperLargeV3,
		FilePath: "./house-speaks-mandarin.mp3",
		Prompt:   "english and mandarin",
	})
	if err != nil {
		return err
	}
	fmt.Println(response.Text)
	return nil
}
