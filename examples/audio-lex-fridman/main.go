// Package main is an example of using groq-go to create a transcription
// using the whisper model.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conneroisu/groq-go"
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
	key := os.Getenv("GROQ_KEY")
	if key == "" {
		return fmt.Errorf("GROQ_KEY is required")
	}
	client, err := groq.NewClient(key)
	if err != nil {
		return err
	}
	response, err := client.Transcribe(ctx, groq.AudioRequest{
		Model:    groq.ModelWhisperLargeV3,
		FilePath: "./The Roman Emperors who went insane Gregory Aldrete and Lex Fridman.mp3",
	})
	if err != nil {
		return err
	}
	fmt.Println(response.Text)
	return nil
}
