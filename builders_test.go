//go:build !test
// +build !test

package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/conneroisu/groq-go/internal/test"
	"github.com/stretchr/testify/assert"
)

// mockFormBuilder is a mock form builder.
type mockFormBuilder struct {
	mockCreateFormFile       func(string, *os.File) error
	mockCreateFormFileReader func(string, io.Reader, string) error
	mockWriteField           func(string, string) error
	mockClose                func() error
}

// CreateFormFile is a mock form builder create form file method.
func (fb *mockFormBuilder) CreateFormFile(
	fieldname string,
	file *os.File,
) error {
	return fb.mockCreateFormFile(fieldname, file)
}

// CreateFormFileReader is a mock form builder create form file reader method
func (fb *mockFormBuilder) CreateFormFileReader(
	fieldname string,
	r io.Reader,
	filename string,
) error {
	return fb.mockCreateFormFileReader(fieldname, r, filename)
}

// WriteField is a mock form builder write field method.
func (fb *mockFormBuilder) WriteField(fieldname, value string) error {
	return fb.mockWriteField(fieldname, value)
}

// Close is a mock form builder close method.
func (fb *mockFormBuilder) Close() error {
	return fb.mockClose()
}

// FormDataContentType is a mock form builder.
func (fb *mockFormBuilder) FormDataContentType() string {
	return ""
}

// failingWriter is a failing writer.
type failingWriter struct{}

var errMockFailingWriterError = errors.New("mock writer failed")

// Write is a failing writer.
func (*failingWriter) Write([]byte) (int, error) {
	return 0, errMockFailingWriterError
}

// TestFormBuilderWithFailingWriter tests the form builder returns an error when the writer fails.
func TestFormBuilderWithFailingWriter(t *testing.T) {
	a := assert.New(t)
	dir, cleanup := test.CreateTestDirectory(t)
	defer cleanup()
	file, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Errorf("Error creating tmp file: %v", err)
	}
	defer file.Close()
	defer os.Remove(file.Name())

	builder := newFormBuilder(&failingWriter{})
	err = builder.CreateFormFile("file", file)
	a.ErrorIs(
		err,
		errMockFailingWriterError,
		"formbuilder should return error if writer fails",
	)
}

// TestFormBuilderWithClosedFile tests the form builder returns an error when the file is closed.
func TestFormBuilderWithClosedFile(t *testing.T) {
	a := assert.New(t)
	dir, cleanup := test.CreateTestDirectory(t)
	defer cleanup()

	file, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Errorf("Error creating tmp file: %v", err)
	}
	file.Close()
	defer os.Remove(file.Name())

	body := &bytes.Buffer{}
	builder := newFormBuilder(body)
	err = builder.CreateFormFile("file", file)
	a.Error(err, "formbuilder should return error if file is closed")
	a.ErrorIs(
		err,
		os.ErrClosed,
		"formbuilder should return error if file is closed",
	)
}

// TestRequestBuilderReturnsRequest  tests the request builder returns a
// request.
func TestRequestBuilderReturnsRequest(t *testing.T) {
	b := newRequestBuilder()
	var (
		ctx         = context.Background()
		method      = http.MethodPost
		url         = "/foo"
		request     = map[string]string{"foo": "bar"}
		reqBytes, _ = json.Marshal(request)
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

// TestRequestBuilderReturnsRequestWhenRequestOfArgsIsNil tests the request
// builder returns a request when the request of args is nil.
func TestRequestBuilderReturnsRequestWhenRequestOfArgsIsNil(t *testing.T) {
	var (
		ctx     = context.Background()
		method  = http.MethodGet
		url     = "/foo"
		want, _ = http.NewRequestWithContext(ctx, method, url, nil)
	)
	b := newRequestBuilder()
	got, _ := b.Build(ctx, method, url, nil, nil)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Build() got = %v, want %v", got, want)
	}
}
