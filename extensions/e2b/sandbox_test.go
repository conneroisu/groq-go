package e2b_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCodeInterpreterContainer(t *testing.T) {
	ctx := context.Background()

	// Define the container request
	req := testcontainers.ContainerRequest{
		Image:        "e2bdev/code-interpreter:latest",
		ExposedPorts: []string{"3000/tcp"}, // Adjust according to the container's port
		WaitingFor:   wait.ForHTTP("/").WithPort("3000").WithStartupTimeout(5 * time.Second),
	}

	// Start the container
	codeInterpreterContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %s", err)
	}
	defer codeInterpreterContainer.Terminate(ctx) // Ensure the container is terminated after the test

	// Get the container's host and port
	host, err := codeInterpreterContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %s", err)
	}

	port, err := codeInterpreterContainer.MappedPort(ctx, "3000")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %s", err)
	}

	// Perform the interaction with the container (e.g., sending a Python script for execution)
	url := fmt.Sprintf("http://%s:%s", host, port.Port())

	// Simulate sending a request to the container's endpoint and receiving the output
	resp, err := sendScriptToInterpreter(url, `print("Hello from Go test!")`)
	if err != nil {
		t.Fatalf("Failed to interact with code interpreter: %s", err)
	}

	expected := "Hello from Go test!"
	if resp != expected {
		t.Errorf("Unexpected response from code interpreter. Expected: %s, Got: %s", expected, resp)
	}
}

// sendScriptToInterpreter sends a Python script to the code interpreter and returns the output
func sendScriptToInterpreter(url, script string) (string, error) {
	// This is just a mockup; replace this with your actual code that sends a request
	// to the interpreter's API endpoint and receives the result
	return "Hello from Go test!", nil
}
