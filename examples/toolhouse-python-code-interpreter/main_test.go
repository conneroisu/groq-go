package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	err := run(ctx)
	a.NoError(err)
}
