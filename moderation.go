package groq

import (
	"context"
	"net/http"
	"strings"

	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/conneroisu/groq-go/pkg/models"
	"github.com/conneroisu/groq-go/pkg/moderation"
)

type (
	// Moderation represents the response of a moderation request.
	Moderation struct {
		// Categories is the categories of the result.
		Categories []moderation.HarmfulCategory `json:"categories"`
		// Flagged is the flagged status of the result.
		Flagged bool `json:"flagged"`
	}
)

// Moderate performs a moderation api call over a string.
// Input can be an array or slice but a string will reduce the complexity.
func (c *Client) Moderate(
	ctx context.Context,
	messages []ChatCompletionMessage,
	model models.ModerationModel,
) (response Moderation, err error) {
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodPost,
		c.fullURL(chatCompletionsSuffix, withModel(model)),
		builders.WithBody(&struct {
			Messages []ChatCompletionMessage `json:"messages"`
			Model    models.ModerationModel  `json:"model,omitempty"`
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
		response.Flagged = true
		split := strings.Split(
			strings.Split(resp.Choices[0].Message.Content, "\n")[1],
			",",
		)
		for _, s := range split {
			response.Categories = append(
				response.Categories,
				moderation.SectionMap[strings.TrimSpace(s)],
			)
		}
	}
	return
}
