package builders

import "net/url"

// WithBody sets the body for a request.
func WithBody(body any) RequestOption {
	return func(args *requestOptions) {
		args.body = body
	}
}

// WithContentType sets the content type for a request.
func WithContentType(contentType string) RequestOption {
	return func(args *requestOptions) {
		args.header.Set("Content-Type", contentType)
	}
}

type (
	// Querier is an interface for a request querier.
	//
	// It allows for modifying the URL before it is sent.
	Querier interface {
		URLQuery(url *url.URL)
	}
)

// WithQuerier sets the querier for a request.
func WithQuerier(querier Querier) RequestOption {
	return func(args *requestOptions) {
		args.querier = querier
	}
}
