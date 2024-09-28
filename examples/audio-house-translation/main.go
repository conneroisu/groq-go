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
	client, err := groq.NewClient(os.Getenv("GROQ_KEY"))
	if err != nil {
		return err
	}
	response, err := client.CreateTranslation(ctx, groq.AudioRequest{
		Model:    groq.ModelWhisperLargeV3,
		FilePath: "./house-speaks-mandarin.mp3",
		Prompt:   "english and mandarin",
	})
	if err != nil {
		return err
	}
	fmt.Println(response.Text)
	return nil
}
