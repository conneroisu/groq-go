package main

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	if len(os.Getenv("UNIT")) < 1 {
		t.Skip("Skipping integration test")
	}
	a := assert.New(t)
	ctx := context.Background()
	err := run(ctx)
	a.NoError(err)
}
