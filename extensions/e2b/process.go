package e2b

import (
	"math/rand"
)

type (
	// Process is a process in the sandbox.
	Process struct {
		ext *Sandbox

		ID  string
		cmd string
		cwd string
		env map[string]string
	}
)

// {"jsonrpc": "2.0", "method": "process_start", "params": ["KkLECSZQiN5B", "cat file0.txt", {"PYTHONUNBUFFERED": "1"}, ""], "id": 12}

// {"jsonrpc": "2.0", "method": "process_start", "params": ["Z9SalhcNx641", "cat file9.txt", {"PYTHONUNBUFFERED": "1"}, ""], "id": 341}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func createProcessID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// StartProcess starts a process in the sandbox.
func (s *Sandbox) StartProcess(
	cmd string,
) (proc Process, err error) {
	proc.ID = createProcessID(12)
	return Process{
		cmd: cmd,
	}, nil
}
