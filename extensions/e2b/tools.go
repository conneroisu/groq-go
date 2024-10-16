package e2b

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/conneroisu/groq-go"
)

type (
	SbFn func(ctx context.Context, s *Sandbox, params *Params) (groq.ChatCompletionMessage, error)
	// ToolingWrapper is a wrapper for groq.Tool that allows for custom functions working with a sandbox.
	ToolingWrapper struct {
		ToolMap map[*groq.Tool]SbFn
	}
	// Params are the parameters for any function call.
	Params struct {
		Path    string `json:"path"`
		Data    string `json:"data"`
		Cmd     string `json:"cmd"`
		Timeout int    `json:"timeout"`
		Cwd     string `json:"cwd"`
		Name    string `json:"name"`
	}
)

// GetTools returns the tools wrapped by the ToolWrapper.
func (t *ToolingWrapper) GetTools() []groq.Tool {
	tools := make([]groq.Tool, 0)
	for tool := range t.ToolMap {
		tools = append(tools, *tool)
	}
	return tools
}

// GetToolFn returns the function for the tool with the
// given name.
func (t *ToolingWrapper) GetToolFn(name string) (SbFn, error) {
	for tool, fn := range t.ToolMap {
		if tool.Function.Name == name {
			return fn, nil
		}
	}
	return nil, fmt.Errorf("tool %s not found", name)
}

var (
	defaultToolWrapper = ToolingWrapper{
		ToolMap: toolMap,
	}
	toolMap = map[*groq.Tool]SbFn{
		&mkdirTool: func(ctx context.Context, s *Sandbox, params *Params) (groq.ChatCompletionMessage, error) {
			err := s.Mkdir(ctx, params.Path)
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			return groq.ChatCompletionMessage{
				Content: fmt.Sprintf("Created directory %s.", params.Path),
				Role:    groq.ChatMessageRoleFunction,
				Name:    "mkdir",
			}, nil
		},
		&lsTool: func(ctx context.Context, s *Sandbox, params *Params) (groq.ChatCompletionMessage, error) {
			res, err := s.Ls(ctx, params.Path)
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			jsonBytes, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			return groq.ChatCompletionMessage{
				Content: string(jsonBytes),
				Role:    groq.ChatMessageRoleFunction,
				Name:    "ls",
			}, nil
		},
		&readTool: func(ctx context.Context, s *Sandbox, params *Params) (groq.ChatCompletionMessage, error) {
			content, err := s.Read(ctx, params.Path)
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			return groq.ChatCompletionMessage{
				Content: string(content),
				Role:    groq.ChatMessageRoleFunction,
				Name:    "read",
			}, nil
		},
		&writeTool: func(ctx context.Context, s *Sandbox, params *Params) (groq.ChatCompletionMessage, error) {
			err := s.Write(ctx, params.Path, []byte(params.Data))
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			return groq.ChatCompletionMessage{
				Content: fmt.Sprintf("Successfully wrote to file %s.", params.Path),
				Role:    groq.ChatMessageRoleFunction,
				Name:    "write",
			}, nil
		},
		&startProcessTool: func(ctx context.Context, s *Sandbox, params *Params) (groq.ChatCompletionMessage, error) {
			proc, err := s.NewProcess(params.Cmd, Process{})
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			events := make(chan Event, 100)
			err = proc.Subscribe(ctx, OnStdout, events)
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			err = proc.Subscribe(ctx, OnStderr, events)
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			err = proc.Start(ctx)
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			buf := new(bytes.Buffer)
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case event := <-events:
						buf.Write([]byte(event.Params.Result.Line))
					case <-proc.Done():
						break
					}
				}
			}()
			<-proc.Done()
			return groq.ChatCompletionMessage{
				Content: buf.String(),
				Role:    groq.ChatMessageRoleFunction,
				Name:    "startProcess",
			}, nil
		},
	}
	mkdirTool = groq.Tool{
		Type: groq.ToolTypeFunction,
		Function: groq.FunctionDefinition{
			Name:        "mkdir",
			Description: "Make a directory in the sandbox file system at a given path",
			Parameters: groq.ParameterDefinition{
				Type: "object",
				Properties: map[string]groq.PropertyDefinition{
					"path": {
						Type:        "string",
						Description: "The path of the directory to create",
					},
				},
				Required:             []string{"path"},
				AdditionalProperties: false,
			},
		},
	}
	lsTool = groq.Tool{
		Type: groq.ToolTypeFunction,
		Function: groq.FunctionDefinition{
			Name:        "ls",
			Description: "List the files and directories in the sandbox file system at a given path",
			Parameters: groq.ParameterDefinition{
				Type: "object",
				Properties: map[string]groq.PropertyDefinition{
					"path": {Type: "string",
						Description: "The path of the directory to list",
					},
				},
				Required:             []string{"path"},
				AdditionalProperties: false,
			},
		},
	}
	readTool = groq.Tool{
		Type: groq.ToolTypeFunction,
		Function: groq.FunctionDefinition{
			Name:        "read",
			Description: "Read the contents of a file in the sandbox file system at a given path",
			Parameters: groq.ParameterDefinition{
				Type: "object",
				Properties: map[string]groq.PropertyDefinition{
					"path": {Type: "string",
						Description: "The path of the file to read",
					},
				},
				Required:             []string{"path"},
				AdditionalProperties: false,
			},
		},
	}
	writeTool = groq.Tool{
		Type: groq.ToolTypeFunction,
		Function: groq.FunctionDefinition{
			Name:        "write",
			Description: "Write to a file in the sandbox file system at a given path",
			Parameters: groq.ParameterDefinition{
				Type: "object",
				Properties: map[string]groq.PropertyDefinition{
					"path": {Type: "string",
						Description: "The relative or absolute path of the file to write to",
					},
					"data": {Type: "string",
						Description: "The data to write to the file",
					},
				},
				Required:             []string{"path", "data"},
				AdditionalProperties: false,
			},
		},
	}
	startProcessTool = groq.Tool{
		Type: groq.ToolTypeFunction,
		Function: groq.FunctionDefinition{
			Name:        "start_process",
			Description: "Start a process in the sandbox.",
			Parameters: groq.ParameterDefinition{
				Type: "object",
				Properties: map[string]groq.PropertyDefinition{
					"cmd": {Type: "string",
						Description: "The command to run to start the process",
					},
					"cwd": {Type: "string",
						Description: "The current working directory of the process",
					},
					"timeout": {Type: "number",
						Description: "The timeout in seconds to run the process",
					},
				},
				Required:             []string{"cmd"},
				AdditionalProperties: false,
			},
		},
	}
)

// RunTooling runs the toolcalls in the response.
func (s *Sandbox) RunTooling(
	ctx context.Context,
	response groq.ChatCompletionResponse,
) ([]groq.ChatCompletionMessage, error) {
	if response.Choices[0].FinishReason != groq.FinishReasonFunctionCall && response.Choices[0].FinishReason != "tool_calls" {
		return nil, fmt.Errorf("not a function call")
	}
	respH := []groq.ChatCompletionMessage{}
	for _, tool := range response.Choices[0].Message.ToolCalls {
		for _, t := range s.toolW.GetTools() {
			if t.Function.Name != tool.Function.Name {
				continue
			}
			resp, err := s.runTool(ctx, t, tool)
			if err != nil {
				return nil, err
			}
			respH = append(respH, resp)
		}
	}
	return respH, nil
}

func (s *Sandbox) runTool(
	ctx context.Context,
	tool groq.Tool,
	call groq.ToolCall,
) (groq.ChatCompletionMessage, error) {
	s.logger.Debug("running tool", "tool", tool.Function.Name, "call", call.Function.Name)
	var params *Params
	err := json.Unmarshal(
		[]byte(call.Function.Arguments),
		&params,
	)
	if err != nil {
		return groq.ChatCompletionMessage{}, err
	}
	fn, err := s.toolW.GetToolFn(tool.Function.Name)
	if err != nil {
		return groq.ChatCompletionMessage{
			Content: fmt.Sprintf("Error running tool (does not exist) %s: %s", tool.Function.Name, err.Error()),
			Role:    groq.ChatMessageRoleFunction,
			Name:    tool.Function.Name,
		}, err
	}
	result, err := fn(ctx, s, params)
	if err != nil {
		return groq.ChatCompletionMessage{
			Content: fmt.Sprintf("Error running tool %s: %s", tool.Function.Name, err.Error()),
			Role:    groq.ChatMessageRoleFunction,
			Name:    tool.Function.Name,
		}, err
	}
	return result, nil
}
