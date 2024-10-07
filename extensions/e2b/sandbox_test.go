package e2b_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/conneroisu/groq-go/extensions/e2b"
	"github.com/stretchr/testify/assert"
)

// func TestCodeInterpreterContainer(t *testing.T) {
//         ctx := context.Background()
//
//         // Define the container request
//         req := testcontainers.ContainerRequest{
//                 Image:        "e2bdev/code-interpreter:latest",
//                 ExposedPorts: []string{"3000/tcp"}, // Adjust according to the container's port
//                 WaitingFor:   wait.ForHTTP("/").WithPort("3000").WithStartupTimeout(5 * time.Second),
//         }
//
//         // Start the container
//         codeInterpreterContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
//                 ContainerRequest: req,
//                 Started:          true,
//         })
//         if err != nil {
//                 t.Fatalf("Failed to start container: %s", err)
//         }
//         defer codeInterpreterContainer.Terminate(ctx) // Ensure the container is terminated after the test
//
//         // Get the container's host and port
//         host, err := codeInterpreterContainer.Host(ctx)
//         if err != nil {
//                 t.Fatalf("Failed to get container host: %s", err)
//         }
//
//         port, err := codeInterpreterContainer.MappedPort(ctx, "3000")
//         if err != nil {
//                 t.Fatalf("Failed to get mapped port: %s", err)
//         }
//
//         // Perform the interaction with the container (e.g., sending a Python script for execution)
//         url := fmt.Sprintf("http://%s:%s", host, port.Port())
//
//         // Simulate sending a request to the container's endpoint and receiving the output
//         resp, err := sendScriptToInterpreter(url, `print("Hello from Go test!")`)
//         if err != nil {
//                 t.Fatalf("Failed to interact with code interpreter: %s", err)
//         }
//
//         expected := "Hello from Go test!"
//         if resp != expected {
//                 t.Errorf("Unexpected response from code interpreter. Expected: %s, Got: %s", expected, resp)
//         }
// }
//
// // sendScriptToInterpreter sends a Python script to the code interpreter and returns the output
// func sendScriptToInterpreter(url, script string) (string, error) {
//         // This is just a mockup; replace this with your actual code that sends a request
//         // to the interpreter's API endpoint and receives the result
//         return "Hello from Go test!", nil
// }

func TestPostSandbox(t *testing.T) {
	a := assert.New(t)
	apiKey := os.Getenv("E2B_API_KEY")
	if apiKey == "" {
		t.Fatal("E2B_API_KEY is not set")
	}
	ll := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	sb, err := e2b.NewSandbox(apiKey, e2b.WithLogger(ll))
	a.NoError(err, "NewSandbox error")
	defer func() {
		err = sb.Close()
		a.NoError(err, "Close error")
	}()
	err = sb.Mkdir("heelo")
	a.NoError(err)
	lsr, err := sb.Ls(".")
	a.NoError(err)
	fmt.Println(lsr)
}

// TestWriteRead tests the Write and Read methods of the Sandbox.
// It creates a new sandbox, writes a file to it, reads the file back, and then closes the sandbox.
// It ensures that the file contents are the same as the original file that was written.
func TestWriteRead(t *testing.T) {
	a := assert.New(t)
	apiKey := os.Getenv("E2B_API_KEY")
	if apiKey == "" {
		t.Fatal("E2B_API_KEY is not set")
	}
	sb, err := e2b.NewSandbox(apiKey)
	a.NoError(err, "NewSandbox error")
	defer func() {
		err = sb.Close()
		a.NoError(err, "Close error")
	}()
	filePath := "test.txt"
	content := "Hello, world!"
	err = sb.Write(filePath, []byte(content))
	a.NoError(err, "Write error")
	readContent, err := sb.Read(filePath)
	a.NoError(err, "Read error")
	a.Equal(content, string(readContent))
	err = sb.Stop()
	a.NoError(err, "Stop error")
}

// ress, err := sb.Ls("/tmp/")
// if err != nil {
//         return Sandbox{}, err
// }
// println(fmt.Sprintf("ress: %v", ress))
// err = sb.Mkdir("/tmp/groq-go")
// if err != nil {
//         return Sandbox{}, err
// }
// // see if it there
// ress, err = sb.Ls("/tmp/")
// if err != nil {
//         return Sandbox{}, err
// }
// println(fmt.Sprintf("ress: %v", ress))
