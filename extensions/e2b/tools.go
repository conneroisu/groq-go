package e2b

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/conneroisu/groq-go"
)

func (s *Sandbox) getTools() []groq.Tool {
	tools := []groq.Tool{
		{
			Type: groq.ToolTypeFunction,
			Function: groq.FunctionDefinition{
				Name:        "mkdir",
				Description: "Make a directory in the sandbox file system at a given path.",
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
			Fn: func(ctx context.Context, args ...string) (groq.ChatCompletionMessage, error) {
				err := s.Mkdir(ctx, args[0])
				if err != nil {
					return groq.ChatCompletionMessage{}, err
				}
				return groq.ChatCompletionMessage{
					Content: args[0],
					Role:    groq.ChatMessageRoleFunction,
					Name:    "mkdir",
				}, nil
			},
		},
		// {
		//         Type: groq.ToolTypeFunction,
		//         Function: groq.FunctionDefinition{
		//                 Name:        "rm",
		//                 Description: "Remove a file or directory in the sandbox file system at a given path.",
		//                 Parameters: groq.ParameterDefinition{
		//                         Type: "object",
		//                         Properties: map[string]groq.PropertyDefinition{
		//                                 "path": {
		//                                         Type:        "string",
		//                                         Description: "The path of the file or directory to remove",
		//                                 },
		//                         },
		//                         Required: []string{
		//                                 "path",
		//                         },
		//                         AdditionalProperties: false,
		//                 },
		//         },
		// },
		{
			Type: groq.ToolTypeFunction,
			Function: groq.FunctionDefinition{
				Name:        "ls",
				Description: "List the files and directories in the sandbox file system at a given path.",
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
			Fn: func(ctx context.Context, args ...string) (groq.ChatCompletionMessage, error) {
				res, err := s.Ls(ctx, args[0])
				if err != nil {
					return groq.ChatCompletionMessage{}, err
				}
				return groq.ChatCompletionMessage{
					Content: res[0].Name,
					Role:    groq.ChatMessageRoleFunction,
					Name:    "ls",
				}, nil
			},
		},
		{
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
			Fn: func(ctx context.Context, args ...string) (groq.ChatCompletionMessage, error) {
				content, err := s.Read(args[0])
				if err != nil {
					return groq.ChatCompletionMessage{}, err
				}
				return groq.ChatCompletionMessage{
					Content: string(content),
					Role:    groq.ChatMessageRoleFunction,
					Name:    "read",
				}, nil
			},
		},
		{
			Type: groq.ToolTypeFunction,
			Function: groq.FunctionDefinition{
				Name:        "write",
				Description: "Write to a file in the sandbox file system at a given path",
				Parameters: groq.ParameterDefinition{
					Type: "object",
					Properties: map[string]groq.PropertyDefinition{
						"path": {Type: "string",
							Description: "The relative or absolute path of the file to write to.",
						},
						"data": {Type: "string",
							Description: "The data to write to the file",
						},
					},
					Required:             []string{"path", "data"},
					AdditionalProperties: false,
				},
			},
			Fn: func(ctx context.Context, args ...string) (groq.ChatCompletionMessage, error) {
				name := args[0]
				data := args[1]
				s.logger.Debug("writing to file", "name", name, "data", data)
				err := s.Write(name, []byte(data))
				if err != nil {
					return groq.ChatCompletionMessage{}, err
				}
				return groq.ChatCompletionMessage{
					Content: "Successfully wrote to file.",
					Role:    groq.ChatMessageRoleFunction,
					Name:    "write",
				}, nil
			},
		},
		{
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
			Fn: func(ctx context.Context, args ...string) (groq.ChatCompletionMessage, error) {
				proc, err := s.NewProcess(args[0], Process{})
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
				for {
					select {
					case <-ctx.Done():
						return groq.ChatCompletionMessage{}, ctx.Err()
					case event := <-events:
						buf.Write([]byte(event.Params.Result.Line))
					case <-proc.Done():
						break
					}
				}
				return groq.ChatCompletionMessage{
					Content: buf.String(),
					Role:    groq.ChatMessageRoleFunction,
					Name:    "startProcess",
				}, nil
			},
		},
	}
	return tools
}

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
		for _, t := range s.getTools() {
			if t.Function.Name != tool.Function.Name {
				continue
			}
			s.logger.Debug("running tool", "tool", t.Function.Name)
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
	params := map[string]any{}
	err := json.Unmarshal(
		[]byte(call.Function.Arguments),
		&params,
	)
	if err != nil {
		return groq.ChatCompletionMessage{}, err
	}
	for _, p := range tool.Function.Parameters.Required {
		if _, ok := params[p]; !ok {
			return groq.ChatCompletionMessage{}, ErrMissingRequiredArgument{
				ToolName: tool.Function.Name,
				ArgName:  p,
			}
		}
	}
	ps := make([]string, 0)
	s.logger.Debug("params", "params", params)
	for k := range tool.Function.Parameters.Properties {
		val, ok := params[k]
		if !ok {
			return groq.ChatCompletionMessage{}, ErrToolArgument{
				ToolName: tool.Function.Name,
				ArgName:  k,
			}
		}
		s.logger.Debug("params", "param", k, "value", val)
		ps = append(ps, fmt.Sprintf("%v", val))
	}
	s.logger.Debug("running tool", "tool", tool.Function.Name, "params", ps)
	result, err := tool.Fn(ctx, ps...)
	if err != nil {
		return groq.ChatCompletionMessage{}, err
	}
	return result, nil
}
