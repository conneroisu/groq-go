package e2b

import (
	"context"

	"github.com/conneroisu/groq-go"
)

func (s *Sandbox) modelMkdir(ctx context.Context, args ...string) (groq.ChatCompletionMessage, error) {
	err := s.Mkdir(args[0])
	if err != nil {
		return groq.ChatCompletionMessage{}, err
	}
	return groq.ChatCompletionMessage{
		Content: args[0],
		Role:    groq.ChatMessageRoleFunction,
		Name:    "mkdir",
	}, nil
}

func (s *Sandbox) modelRead(ctx context.Context, args ...string) (groq.ChatCompletionMessage, error) {
	content, err := s.Read(args[0])
	if err != nil {
		return groq.ChatCompletionMessage{}, err
	}
	return groq.ChatCompletionMessage{
		Content: string(content),
		Role:    groq.ChatMessageRoleFunction,
		Name:    "read",
	}, nil
}

func (s *Sandbox) modelWrite(ctx context.Context, args ...string) (groq.ChatCompletionMessage, error) {
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
}

func (s *Sandbox) modelLs(ctx context.Context, args ...string) (groq.ChatCompletionMessage, error) {
	res, err := s.Ls(args[0])
	if err != nil {
		return groq.ChatCompletionMessage{}, err
	}
	return groq.ChatCompletionMessage{
		Content: res[0].Name,
		Role:    groq.ChatMessageRoleFunction,
		Name:    "ls",
	}, nil
}

func (s *Sandbox) modelStartProcess(ctx context.Context, args ...string) (groq.ChatCompletionMessage, error) {
	proc, err := s.NewProcess(args[0], &Process{})
	if err != nil {
		return groq.ChatCompletionMessage{}, err
	}
	return groq.ChatCompletionMessage{
		Content: proc.id,
		Role:    groq.ChatMessageRoleFunction,
		Name:    "startProcess",
	}, nil
}
