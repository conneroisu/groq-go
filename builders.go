package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

// FormBuilder is an interface for building a form.
type FormBuilder interface {
	CreateFormFile(fieldname string, file *os.File) error
	CreateFormFileReader(fieldname string, r io.Reader, filename string) error
	WriteField(fieldname, value string) error
	Close() error
	FormDataContentType() string
}

// DefaultFormBuilder is a default implementation of FormBuilder.
type DefaultFormBuilder struct {
	writer *multipart.Writer
}

// NewFormBuilder creates a new DefaultFormBuilder.
func NewFormBuilder(body io.Writer) *DefaultFormBuilder {
	return &DefaultFormBuilder{
		writer: multipart.NewWriter(body),
	}
}

// CreateFormFile creates a form file.
func (fb *DefaultFormBuilder) CreateFormFile(
	fieldname string,
	file *os.File,
) error {
	return fb.createFormFile(fieldname, file, file.Name())
}

// CreateFormFileReader creates a form file from a reader.
func (fb *DefaultFormBuilder) CreateFormFileReader(
	fieldname string,
	r io.Reader,
	filename string,
) error {
	return fb.createFormFile(fieldname, r, path.Base(filename))
}

// createFormFile creates a form file.
func (fb *DefaultFormBuilder) createFormFile(
	fieldname string,
	r io.Reader,
	filename string,
) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	fieldWriter, err := fb.writer.CreateFormFile(fieldname, filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(fieldWriter, r)
	if err != nil {
		return err
	}

	return nil
}

// WriteField writes a field to the form.
func (fb *DefaultFormBuilder) WriteField(fieldname, value string) error {
	return fb.writer.WriteField(fieldname, value)
}

// Close closes the form.
func (fb *DefaultFormBuilder) Close() error {
	return fb.writer.Close()
}

// FormDataContentType returns the content type of the form.
func (fb *DefaultFormBuilder) FormDataContentType() string {
	return fb.writer.FormDataContentType()
}

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
