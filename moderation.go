package groq

import (
	"context"
	"net/http"
	"strings"

	"github.com/conneroisu/groq-go/pkg/builders"
)

// Moderate performs a moderation api call over a string.
// Input can be an array or slice but a string will reduce the complexity.
func (c *Client) Moderate(
	ctx context.Context,
	messages []ChatCompletionMessage,
	model ModerationModel,
) (response []Moderation, err error) {
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodPost,
		c.fullURL(chatCompletionsSuffix, withModel(model)),
		builders.WithBody(&struct {
			Messages []ChatCompletionMessage `json:"messages"`
			Model    ModerationModel         `json:"model,omitempty"`
		}{
			Messages: messages,
			Model:    model,
		}),
	)
	if err != nil {
		return
	}
	var resp ChatCompletionResponse
	err = c.sendRequest(req, &resp)
	if err != nil {
		return
	}
	if strings.Contains(resp.Choices[0].Message.Content, "unsafe") {
		split := strings.Split(
			strings.Split(resp.Choices[0].Message.Content, "\n")[1],
			",",
		)
		for _, s := range split {
			response = append(
				response,
				SectionMap[strings.TrimSpace(s)],
			)
		}
	}
	return
}
