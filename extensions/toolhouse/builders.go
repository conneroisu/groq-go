package toolhouse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// formBuilder is an interface for building a form.
type (
	requestBuilder interface {
		Build(
			ctx context.Context,
			method, url string,
			body any,
			header http.Header,
		) (*http.Request, error)
	}
	httpRequestBuilder struct{}
	requestOptions     struct {
		body   any
		header http.Header
	}
	requestOption func(*requestOptions)
)

func (e *Extension) setCommonHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Toolhouse/1.2.1 Python/3.11.0")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.apiKey))
	req.Header.Set("Content-Type", applicationJSON)
}

func newRequestBuilder() requestBuilder {
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
