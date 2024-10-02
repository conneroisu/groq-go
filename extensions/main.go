package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/conneroisu/groq-go/extensions/toolhouse"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// run runs the main function.
func run() error {
	ctx := context.Background()
	ext, err := toolhouse.NewExtension(os.Getenv("TOOLHOUSE_API_KEY"))
	if err != nil {
		return err
	}
	resp, err := ext.GetTools(
		ctx,
		toolhouse.WithProvider("openai"),
		toolhouse.WithMetadata(
			map[string]any{
				"id":       "conner",
				"timezone": 5,
			},
		),
	)
	if err != nil {
		return err
	}
	jsV, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsV))
	return nil
}
