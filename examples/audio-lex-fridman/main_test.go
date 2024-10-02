package main

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestAudioLexFridman(t *testing.T) {
	if len(os.Getenv("UNIT")) < 1 {
		t.Skip("Skipping AudioLexFridman test")
	}
	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
