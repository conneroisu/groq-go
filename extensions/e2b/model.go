package e2b

import (
	"context"
	"io"
	"time"
)

type (
	// Requester is an interface for an instance that sends rpc requests.
	//
	// Implementations should be conccurent safe.
	Requester interface {
		Write(
			ctx context.Context,
			method Method,
			params []any,
			respCh chan []byte,
		)
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
	}
	// Processor is an interface for a process.
	Processor interface {
		Start(
			ctx context.Context,
			cmd string,
			timeout time.Duration,
		)
		Subscriber
	}
	// Lifer is an interface for keeping sandboxes alive.
	Lifer interface {
		// KeepAlive keeps the underlying interface alive.
		//
		// If the context is cancelled before requesting the timeout,
		// the error will be ctx.Err().
		KeepAlive(ctx context.Context, timeout time.Duration) error
	}
	// Subscriber is an interface for an instance that can subscribe to an event.
	Subscriber interface {
		Subscribe(
			ctx context.Context,
			event ProcessEvents,
			eCh chan<- Event,
		)
	}
	// Watcher is an interface for a instance that can watch a filesystem.
	Watcher interface {
	}
)
