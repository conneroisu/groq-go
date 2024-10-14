package e2b_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/conneroisu/groq-go/extensions/e2b"
	"github.com/stretchr/testify/assert"
)

var (
	defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{}
			}
			if a.Key == "level" {
				return slog.Attr{}
			}
			if a.Key == "source" {
				str := a.Value.String()
				split := strings.Split(str, "/")
				if len(split) > 2 {
					a.Value = slog.StringValue(strings.Join(split[len(split)-2:], "/"))
				}
			}
			return a
		}}))
)

func getapiKey(t *testing.T) string {
	apiKey := os.Getenv("E2B_API_KEY")
	if apiKey == "" {
		t.Fail()
	}
	return apiKey
}

func TestPostSandbox(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	sb, err := e2b.NewSandbox(
		ctx,
		getapiKey(t),
		e2b.WithLogger(defaultLogger),
	)
	a.NoError(err, "NewSandbox error")
	lsr, err := sb.Ls(ctx, ".")
	a.NoError(err)
	for _, name := range []string{"boot", "code", "dev", "etc", "home"} {
		a.Contains(lsr, e2b.LsResult{
			Name:  name,
			IsDir: true,
		})
	}
	err = sb.Mkdir(ctx, "heelo")
	a.NoError(err)
	lsr, err = sb.Ls(ctx, "/")
	a.NoError(err)
	a.Contains(lsr, e2b.LsResult{
		Name:  "heelo",
		IsDir: true,
	})
}

// TestWriteRead tests the Write and Read methods of the Sandbox.
func TestWriteRead(t *testing.T) {
	filePath := "test.txt"
	content := "Hello, world!"
	a := assert.New(t)
	ctx := context.Background()
	sb, err := e2b.NewSandbox(
		ctx,
		getapiKey(t),
		e2b.WithLogger(defaultLogger),
	)
	a.NoError(err, "NewSandbox error")
	err = sb.Write(filePath, []byte(content))
	a.NoError(err, "Write error")
	readContent, err := sb.Read(filePath)
	a.NoError(err, "Read error")
	a.Equal(content, string(readContent), "Read content does not match written content")
	readBytesContent, err := sb.ReadBytes(ctx, filePath)
	a.NoError(err, "ReadBytes error")
	a.Equal(content, string(readBytesContent), "ReadBytes content does not match written content")
	err = sb.Stop(ctx)
	a.NoError(err, "Stop error")
}

func TestCreateProcess(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	sb, err := e2b.NewSandbox(
		ctx,
		getapiKey(t),
		e2b.WithLogger(defaultLogger),
	)
	a.NoError(err, "NewSandbox error")
	proc, err := sb.NewProcess("echo 'Hello World!'",
		e2b.Process{
			Env: map[string]string{
				"FOO": "bar",
			},
		})
	a.NoError(err, "could not create process")
	err = proc.Start()
	a.NoError(err)
	proc, err = sb.NewProcess("sleep 2 && echo 'Hello World!'", e2b.Process{})
	a.NoError(err, "could not create process")
	err = proc.Start()
	a.NoError(err)
	ctx, cancel := context.WithTimeout(ctx, time.Second*6)
	defer cancel()
	events := make(chan e2b.Event, 10)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-events:
				jsonBytes, err := json.MarshalIndent(event, "", "  ")
				if err != nil {
					a.Error(err)
					return
				}
				print(string(jsonBytes))
			}
		}
	}()
	err = proc.Subscribe(ctx, e2b.OnStdout, events)
	a.NoError(err)
	time.Sleep(3 * time.Second)
}

func TestFilesystemSubscribe(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	sb, err := e2b.NewSandbox(
		ctx,
		getapiKey(t),
		e2b.WithLogger(defaultLogger),
		e2b.WithCwd("/tmp"),
	)
	a.NoError(err, "NewSandbox error")
	// subscribe to a file
	events := make(chan e2b.Event)
	err = sb.Watch(ctx, "/tmp/", events)
	a.NoError(err)
	go func() {
		for event := range events {
			jsonBytes, err := json.MarshalIndent(event, "", "  ")
			if err != nil {
				a.Error(err)
				return
			}
			print(string(jsonBytes))
		}
	}()
	// create a file
	err = sb.Write("/tmp/file.txt", []byte("Hello World!"))
	a.NoError(err)
	err = sb.Write("/tmp/file2.txt", []byte("Hello World!"))
	a.NoError(err)
	time.Sleep(3 * time.Second)
}

func TestKeepAlive(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	sb, err := e2b.NewSandbox(
		ctx,
		getapiKey(t),
		e2b.WithLogger(defaultLogger),
	)
	a.NoError(err, "NewSandbox error")
	err = sb.KeepAlive(time.Minute * 2)
	a.NoError(err)
}
