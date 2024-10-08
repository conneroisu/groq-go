package e2b

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type (
	requestOption func(*requestOptions)

	requestOptions struct {
		body   any
		header http.Header
	}
	requestBuilder interface {
		Build(
			ctx context.Context,
			method, url string,
			body any,
			header http.Header,
		) (*http.Request, error)
	}
	httpRequestBuilder struct{}
)

func (s *Sandbox) setCommonHeaders(req *http.Request) {
	req.Header.Set("X-API-Key", s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}

func (s *Sandbox) newRequest(
	ctx context.Context,
	method, url string,
	setters ...requestOption,
) (*http.Request, error) {
	// Default Options
	args := &requestOptions{
		body:   nil,
		header: http.Header{},
	}
	for _, setter := range setters {
		setter(args)
	}
	req, err := s.requestBuilder.Build(
		ctx,
		method,
		url,
		args.body,
		args.header,
	)
	if err != nil {
		return nil, err
	}
	s.setCommonHeaders(req)
	return req, nil
}

func newRequestBuilder() *httpRequestBuilder {
	return &httpRequestBuilder{}
}

func (b *httpRequestBuilder) Build(
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
	req, err = http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return
	}
	if header != nil {
		req.Header = header
	}
	return
}

func withBody(body any) requestOption {
	return func(args *requestOptions) {
		args.body = body
	}
}
