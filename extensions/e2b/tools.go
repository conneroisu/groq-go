package e2b

import (
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
					Required: []string{
						"path",
					},
					AdditionalProperties: false,
				},
			},
			Fn: s.modelMkdir,
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
						"path": {
							Type:        "string",
							Description: "The path of the directory to list",
						},
					},
					Required: []string{
						"path",
					},
					AdditionalProperties: false,
				},
			},
			Fn: s.modelLs,
		},
		{
			Type: groq.ToolTypeFunction,
			Function: groq.FunctionDefinition{
				Name:        "read",
				Description: "Read the contents of a file in the sandbox file system at a given path",
				Parameters: groq.ParameterDefinition{
					Type: "object",
					Properties: map[string]groq.PropertyDefinition{
						"path": {
							Type:        "string",
							Description: "The path of the file to read",
						},
					},
					Required: []string{
						"path",
					},
					AdditionalProperties: false,
				},
			},
			Fn: s.modelRead,
		},
		{
			Type: groq.ToolTypeFunction,
			Function: groq.FunctionDefinition{
				Name:        "write",
				Description: "Write to a file in the sandbox file system at a given path",
				Parameters: groq.ParameterDefinition{
					Type: "object",
					Properties: map[string]groq.PropertyDefinition{
						"path": {
							Type:        "string",
							Description: "The relative or absolute path of the file to write to.",
						},
						"data": {
							Type:        "string",
							Description: "The data to write to the file",
						},
					},
					Required: []string{
						"path",
						"data",
					},
					AdditionalProperties: false,
				},
			},
			Fn: s.modelWrite,
		},
		{
			Type: groq.ToolTypeFunction,
			Function: groq.FunctionDefinition{
				Name:        "start_process",
				Description: "Start a process in the sandbox.",
				Parameters: groq.ParameterDefinition{
					Type: "object",
					Properties: map[string]groq.PropertyDefinition{
						"cmd": {
							Type:        "string",
							Description: "The command to run to start the process",
						},
						"cwd": {
							Type:        "string",
							Description: "The current working directory of the process",
						},
						"timeout": {
							Type:        "number",
							Description: "The timeout in seconds to run the process",
						},
					},
					Required: []string{
						"cmd",
					},
					AdditionalProperties: false,
				},
			},
			Fn: s.modelStartProcess,
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
