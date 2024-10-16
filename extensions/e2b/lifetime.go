package e2b

import (
	"context"
	"time"
)

type (
	// Lifer is an interface for keeping sandboxes alive.
	Lifer interface {
		// KeepAlive keeps the underlying interface alive.
		//
		// If the context is cancelled before requesting the timeout,
		// the error will be ctx.Err().
		KeepAlive(ctx context.Context, timeout time.Duration) error
	}
)
