package play

import "net/http"

type Option func(*PlayAI)

func WithBaseURL(baseURL string) Option {
	return func(p *PlayAI) { p.baseURL = baseURL }
}

func WithClient(client *http.Client) Option {
	return func(p *PlayAI) { p.client = client }
}
