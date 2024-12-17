package e2b

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/tools"
)

type (
	// SbFn is a function that can be used to run a tool.
	SbFn func(ctx context.Context, s *Sandbox, params *Params) (groq.ChatCompletionMessage, error)
	// ToolingWrapper is a wrapper for tools.Tool that allows for custom functions working with a sandbox.
	ToolingWrapper struct {
		ToolMap map[*tools.Tool]SbFn
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

// getTools returns the tools wrapped by the ToolWrapper.
func (t *ToolingWrapper) getTools() []tools.Tool {
	tools := make([]tools.Tool, 0)
	for tool := range t.ToolMap {
		tools = append(tools, *tool)
	}
	return tools
}

// GetTools returns the tools wrapped by the ToolWrapper.
func (s *Sandbox) GetTools() []tools.Tool {
	return s.toolW.getTools()
}

// GetToolFn returns the function for the tool with the
// given name.
func (t *ToolingWrapper) GetToolFn(name string) (SbFn, error) {
	for tool, fn := range t.ToolMap {
		if tool.Function.Name == name {
			return fn, nil
		}
	}
	return nil, fmt.Errorf("Error running tool (does not exist) %s", name)
}

var (
	defaultToolWrapper = ToolingWrapper{
		ToolMap: toolMap,
	}
	toolMap = map[*tools.Tool]SbFn{
		&mkdirTool: func(
			ctx context.Context,
			s *Sandbox,
			params *Params,
		) (groq.ChatCompletionMessage, error) {
			err := s.Mkdir(ctx, params.Path)
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			return groq.ChatCompletionMessage{
				Content: fmt.Sprintf("Created directory %s.", params.Path),
				Role:    groq.RoleFunction,
				Name:    "mkdir",
			}, nil
		},
		&lsTool: func(
			ctx context.Context,
			s *Sandbox,
			params *Params,
		) (groq.ChatCompletionMessage, error) {
			res, err := s.Ls(ctx, params.Path)
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			jsonBytes, err := json.Marshal(res)
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			return groq.ChatCompletionMessage{
				Content: string(jsonBytes),
				Role:    groq.RoleFunction,
				Name:    "ls",
			}, nil
		},
		&readTool: func(
			ctx context.Context,
			s *Sandbox,
			params *Params,
		) (groq.ChatCompletionMessage, error) {
			content, err := s.Read(ctx, params.Path)
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			return groq.ChatCompletionMessage{
				Content: string(content),
				Role:    groq.RoleFunction,
				Name:    "read",
			}, nil
		},
		&writeTool: func(
			ctx context.Context,
			s *Sandbox,
			params *Params,
		) (groq.ChatCompletionMessage, error) {
			err := s.Write(ctx, params.Path, []byte(params.Data))
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			return groq.ChatCompletionMessage{
				Content: fmt.Sprintf("Successfully wrote to file %s.", params.Path),
				Role:    groq.RoleFunction,
				Name:    "write",
			}, nil
		},
		&startProcessTool: func(
			ctx context.Context,
			s *Sandbox,
			params *Params,
		) (groq.ChatCompletionMessage, error) {
			proc, err := s.NewProcess(params.Cmd)
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			e, errCh := proc.SubscribeStdout(ctx)
			if err != nil {
				return groq.ChatCompletionMessage{}, err
			}
			e2, errCh := proc.SubscribeStderr(ctx)
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
					case event := <-e:
						buf.Write([]byte(event.Params.Result.Line))
						continue
					case event := <-e2:
						buf.Write([]byte(event.Params.Result.Line))
						continue
					case <-errCh:
						return
					case <-proc.Done():
						return
					}
				}
			}()
			<-proc.Done()
			return groq.ChatCompletionMessage{
				Content: buf.String(),
				Role:    groq.RoleFunction,
				Name:    "startProcess",
			}, nil
		},
	}
	mkdirTool = tools.Tool{
		Type: tools.ToolTypeFunction,
		Function: tools.FunctionDefinition{
			Name:        "mkdir",
			Description: "Make a directory in the sandbox file system at a given path",
			Parameters: tools.FunctionParameters{
				Type: "object",
				Properties: map[string]tools.PropertyDefinition{
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
	lsTool = tools.Tool{
		Type: tools.ToolTypeFunction,
		Function: tools.FunctionDefinition{
			Name:        "ls",
			Description: "List the files and directories in the sandbox file system at a given path",
			Parameters: tools.FunctionParameters{
				Type: "object",
				Properties: map[string]tools.PropertyDefinition{
					"path": {Type: "string",
						Description: "The path of the directory to list",
					},
				},
				Required:             []string{"path"},
				AdditionalProperties: false,
			},
		},
	}
	readTool = tools.Tool{
		Type: tools.ToolTypeFunction,
		Function: tools.FunctionDefinition{
			Name:        "read",
			Description: "Read the contents of a file in the sandbox file system at a given path",
			Parameters: tools.FunctionParameters{
				Type: "object",
				Properties: map[string]tools.PropertyDefinition{
					"path": {Type: "string",
						Description: "The path of the file to read",
					},
				},
				Required:             []string{"path"},
				AdditionalProperties: false,
			},
		},
	}
	writeTool = tools.Tool{
		Type: tools.ToolTypeFunction,
		Function: tools.FunctionDefinition{
			Name:        "write",
			Description: "Write to a file in the sandbox file system at a given path",
			Parameters: tools.FunctionParameters{
				Type: "object",
				Properties: map[string]tools.PropertyDefinition{
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
	startProcessTool = tools.Tool{
		Type: tools.ToolTypeFunction,
		Function: tools.FunctionDefinition{
			Name:        "start_process",
			Description: "Start a process in the sandbox.",
			Parameters: tools.FunctionParameters{
				Type: "object",
				Properties: map[string]tools.PropertyDefinition{
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
	if response.Choices[0].FinishReason != groq.ReasonFunctionCall && response.Choices[0].FinishReason != "tool_calls" {
		return nil, fmt.Errorf("not a function call: %v", response.Choices[0].FinishReason)
	}
	respH := []groq.ChatCompletionMessage{}
	for _, tool := range response.Choices[0].Message.ToolCalls {
		for _, t := range s.toolW.getTools() {
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
	tool tools.Tool,
	call tools.ToolCall,
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
			Content: err.Error(),
			Role:    groq.RoleFunction,
			Name:    tool.Function.Name,
		}, err
	}
	result, err := fn(ctx, s, params)
	if err != nil {
		return groq.ChatCompletionMessage{
			Content: fmt.Sprintf("Error running tool %s: %s", tool.Function.Name, err.Error()),
			Role:    groq.RoleFunction,
			Name:    tool.Function.Name,
		}, err
	}
	return result, nil
}
