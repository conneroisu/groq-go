package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// RequestBuilder is an interface that defines the Build method.
type RequestBuilder interface {
	Build(
		ctx context.Context,
		method, url string,
		body any,
		header http.Header,
	) (*http.Request, error)
}

// HTTPRequestBuilder is a struct that implements the RequestBuilder interface.
type HTTPRequestBuilder struct {
	marshaller Marshaller
}

// Marshaller is an interface that defines the Marshal method.
type Marshaller interface {
	Marshal(v any) ([]byte, error)
}

// JSONMarshaller is a struct that implements the Marshaller interface.
type JSONMarshaller struct{}

// Marshal marshals the given value to JSON.
func (j *JSONMarshaller) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// NewRequestBuilder returns a new HTTPRequestBuilder.
func NewRequestBuilder() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{
		marshaller: &JSONMarshaller{},
	}
}

// Build builds a new request.
func (b *HTTPRequestBuilder) Build(
	ctx context.Context,
	method string,
	url string,
	body any,
	header http.Header,
) (req *http.Request, err error) {
	var bodyReader io.Reader
	if body != nil {
		if v, ok := body.(io.Reader); ok {
			bodyReader = v
		} else {
			var reqBytes []byte
			reqBytes, err = b.marshaller.Marshal(body)
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
