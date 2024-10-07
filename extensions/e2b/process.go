package e2b

type (
	// Process is a process in the sandbox.
	Process struct {
		ext *Sandbox

		cmd string
	}
)
