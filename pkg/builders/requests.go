package builders

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

var builder = NewRequestBuilder()

type (
	// Header is an struct interface for setting common headers.
	Header struct {
		SetCommonHeaders func(req *http.Request)
	}
	// RequestBuilder is an interface for building requests.
	RequestBuilder interface {
		Build(
			ctx context.Context,
			method, url string,
			body any,
			header http.Header,
		) (*http.Request, error)
	}
	defaultRequestBuilder struct{}
	requestOptions        struct {
		body   any
		header http.Header
	}
	// RequestOption is an option for a request.
	RequestOption func(*requestOptions)
)

// NewRequestBuilder creates a new default request builder.
func NewRequestBuilder() RequestBuilder {
	return &defaultRequestBuilder{}
}

func (b *defaultRequestBuilder) Build(
	ctx context.Context,
	method string,
	url string,
	body any,
	header http.Header,
) (req *http.Request, err error) {
	var bodyReader io.Reader
	if body != nil {
		v, ok := body.(io.Reader)
		if ok {
			bodyReader = v
		} else {
			var reqBytes []byte
			reqBytes, err = json.Marshal(body)
			if err != nil {
				return
			}
			bodyReader = bytes.NewBuffer(reqBytes)
		}
	}
	req, err = http.NewRequestWithContext(
		ctx,
		method,
		url,
		bodyReader,
	)
	if err != nil {
		return
	}
	if header != nil {
		req.Header = header
	}
	return
}

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

// NewRequest creates a new request.
func NewRequest(
	ctx context.Context,
	c Header,
	method, url string,
	setters ...RequestOption,
) (*http.Request, error) {
	args := &requestOptions{
		body:   nil,
		header: http.Header{},
	}
	for _, setter := range setters {
		setter(args)
	}
	req, err := builder.Build(
		ctx,
		method,
		url,
		args.body,
		args.header,
	)
	if err != nil {
		return nil, err
	}
	c.SetCommonHeaders(req)
	return req, nil
}
