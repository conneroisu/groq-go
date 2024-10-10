package builders

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type (
	Requester interface {
		setCommonHeaders(req *http.Request)
		RequestBuilder
	}
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
	RequestOption func(*requestOptions)
)

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

// NewRequest creates a new request.
func NewRequest(
	ctx context.Context,
	c Requester,
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
	req, err := c.Build(
		ctx,
		method,
		url,
		args.body,
		args.header,
	)
	if err != nil {
		return nil, err
	}
	c.setCommonHeaders(req)
	return req, nil
}
