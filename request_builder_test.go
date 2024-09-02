package groq //nolint:testpackage // testing private field

import (
	"bytes"
	"context"
	"net/http"
	"reflect"
	"testing"
)

// TestRequestBuilderReturnsRequest  tests the request builder returns a request.
func TestRequestBuilderReturnsRequest(t *testing.T) {
	b := NewRequestBuilder()
	var (
		ctx         = context.Background()
		method      = http.MethodPost
		url         = "/foo"
		request     = map[string]string{"foo": "bar"}
		reqBytes, _ = b.marshaller.Marshal(request)
		want, _     = http.NewRequestWithContext(
			ctx,
			method,
			url,
			bytes.NewBuffer(reqBytes),
		)
	)
	got, _ := b.Build(ctx, method, url, request, nil)
	if !reflect.DeepEqual(got.Body, want.Body) ||
		!reflect.DeepEqual(got.URL, want.URL) ||
		!reflect.DeepEqual(got.Method, want.Method) {
		t.Errorf("Build() got = %v, want %v", got, want)
	}
}

// TestRequestBuilderReturnsRequestWhenRequestOfArgsIsNil tests the request builder returns a request when the request of args is nil.
func TestRequestBuilderReturnsRequestWhenRequestOfArgsIsNil(t *testing.T) {
	var (
		ctx     = context.Background()
		method  = http.MethodGet
		url     = "/foo"
		want, _ = http.NewRequestWithContext(ctx, method, url, nil)
	)
	b := NewRequestBuilder()
	got, _ := b.Build(ctx, method, url, nil, nil)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Build() got = %v, want %v", got, want)
	}
}
