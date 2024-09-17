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
	client, err := groq.NewClient(os.Getenv("GROQ_KEY"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	response, err := client.CreateTranscription(context.Background(), groq.AudioRequest{
		Model:    groq.WhisperLargeV3,
		FilePath: "./The Roman Emperors who went insane Gregory Aldrete and Lex Fridman.mp3",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(response.Text)
}
