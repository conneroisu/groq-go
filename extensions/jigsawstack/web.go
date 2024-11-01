package jigsawstack

import (
	"context"
	"net/http"
	"time"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	webSearchEndpoint Endpoint = "v1/web/search"
)

type (
	// WebSearchSuggestions is the response for the web search suggestions
	// api.
	WebSearchSuggestions struct {
		Success     bool     `json:"success"`
		Suggestions []string `json:"suggestions"`
	}
	// WebSearchResponse is the response for the web search api.
	WebSearchResponse struct {
		Success    bool   `json:"success"`
		Query      string `json:"query"`
		SpellFixed string `json:"spell_fixed"`
		IsSafe     bool   `json:"is_safe"`
		AiOverview string `json:"ai_overview"`
		Results    []struct {
			Title        string    `json:"title"`
			URL          string    `json:"url"`
			Description  string    `json:"description"`
			Content      string    `json:"content"`
			SiteName     string    `json:"site_name"`
			SiteLongName string    `json:"site_long_name"`
			Age          time.Time `json:"age"`
			Language     string    `json:"language"`
			IsSafe       bool      `json:"is_safe"`
			Favicon      string    `json:"favicon"`
			Snippets     []string  `json:"snippets"`
			RelatedIndex []struct {
				Title       string `json:"title"`
				URL         string `json:"url"`
				Description string `json:"description"`
				IsSafe      bool   `json:"is_safe"`
			} `json:"related_index,omitempty"`
			Thumbnail string `json:"thumbnail,omitempty"`
		} `json:"results"`
	}
)

// WebSearch performs a web search api call over a query string.
//
// POST https://api.jigsawstack.com/v1/web/search
//
// https://docs.jigsawstack.com/api-reference/web/search
func (j *JigsawStack) WebSearch(
	ctx context.Context,
	query string,
) (response WebSearchResponse, err error) {
	// TODO: may need to santize the query
	uri := j.baseURL + string(webSearchEndpoint) + "?query=" + query
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		uri,
	)
	if err != nil {
		return
	}
	var resp WebSearchResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}

// WebSearchSuggestions performs a web search suggestions api call over a query
// string.
//
// POST https://api.jigsawstack.com/v1/web/search
//
// https://docs.jigsawstack.com/api-reference/web/search
func (j *JigsawStack) WebSearchSuggestions(
	ctx context.Context,
	query string,
) (response WebSearchSuggestions, err error) {
	// TODO: may need to santize the query
	uri := j.baseURL + string(webSearchEndpoint) + "?query=" + query
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		uri,
	)
	if err != nil {
		return
	}
	var resp WebSearchSuggestions
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}
