package e2b

import (
	"context"
	"io"
	"time"
)

type (
	// Receiver is an interface for a constantly receiving instance that
	// can closed.
	//
	// Implementations should be conccurent safe.
	Receiver interface {
		Read(ctx context.Context) error
		io.Closer
	}
	// Identifier is an interface for a constantly running process to
	// identify new request ids.
	Identifier interface {
		Identify(ctx context.Context)
	}
	// Sandboxer is an interface for a sandbox.
	Sandboxer interface {
		// KeepAlive keeps the underlying interface alive.
		//
		// If the context is cancelled before requesting the timeout,
		// the error will be ctx.Err().
		KeepAlive(
			ctx context.Context,
			timeout time.Duration,
		) error
		// NewProcess creates a new process.
		NewProcess(
			cmd string,
		) (*Processor, error)

		// Write writes a file to the filesystem.
		Write(
			ctx context.Context,
			method Method,
			params []any,
			respCh chan<- []byte,
		)
		// Read reads a file from the filesystem.
		Read(
			ctx context.Context,
			path string,
		) (string, error)
	}
	// Processor is an interface for a process.
	Processor interface {
		Start(
			ctx context.Context,
			cmd string,
			timeout time.Duration,
		)
		SubscribeStdout() (events chan Event, err error)
		SubscribeStderr() (events chan Event, err error)
	}
	// Watcher is an interface for a instance that can watch a filesystem.
	Watcher interface {
		Watch(
			ctx context.Context,
			path string,
		) (<-chan Event, error)
	}
)
