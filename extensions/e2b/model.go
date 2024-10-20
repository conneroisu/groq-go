package e2b

import (
	"context"
	"io"
	"time"
)

type (
	// Requester is an interface for an instance that sends requests to a filesystem.
	//
	// Implementations should be conccurent safe.
	Requester interface {
		// Write writes a file to the filesystem.
		Write(
			ctx context.Context,
			method Method,
			params []any,
			respCh chan []byte,
		)
		// Read reads a file from the filesystem.
		Read(
			ctx context.Context,
			path string,
		) (string, error)
	}
	// Receiver is an interface for a constantly receiving instance.
	//
	// Implementations should be conccurent safe.
	Receiver interface {
		Read(ctx context.Context) error
		io.Closer
	}
	// Identifier is an interface for a constantly running process to identify new request ids.
	Identifier interface {
		Identify(ctx context.Context)
	}
	// Sandboxer is an interface for a sandbox.
	Sandboxer interface {
		Lifer
		// NewProcess creates a new process.
		NewProcess(
			cmd string,
		) (*Processor, error)
	}
	// Processor is an interface for a process.
	Processor interface {
		Start(
			ctx context.Context,
			cmd string,
			timeout time.Duration,
		)
		Subscribe(
			ctx context.Context,
			event ProcessEvents,
			eCh chan<- Event,
		)
	}
	// Lifer is an interface for keeping sandboxes alive.
	Lifer interface {
		// KeepAlive keeps the underlying interface alive.
		//
		// If the context is cancelled before requesting the timeout,
		// the error will be ctx.Err().
		KeepAlive(ctx context.Context, timeout time.Duration) error
	}
	// Watcher is an interface for a instance that can watch a filesystem.
	Watcher interface {
		Watch(
			ctx context.Context,
			path string,
		) (<-chan Event, error)
	}
)
