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

// formBuilder is an interface for building a form.
type (
	formBuilder interface {
		io.Closer
		CreateFormFile(fieldname string, file *os.File) error
		CreateFormFileReader(fieldname string, r io.Reader, filename string) error
		WriteField(fieldname, value string) error
		FormDataContentType() string
	}
	defaultFormBuilder struct {
		writer *multipart.Writer
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

// newFormBuilder creates a new DefaultFormBuilder.
func newFormBuilder(body io.Writer) *defaultFormBuilder {
	return &defaultFormBuilder{
		writer: multipart.NewWriter(body),
	}
}

func (fb *defaultFormBuilder) CreateFormFile(
	fieldname string,
	file *os.File,
) error {
	return fb.createFormFile(fieldname, file, file.Name())
}

func (fb *defaultFormBuilder) CreateFormFileReader(
	fieldname string,
	r io.Reader,
	filename string,
) error {
	return fb.createFormFile(fieldname, r, path.Base(filename))
}

func (fb *defaultFormBuilder) createFormFile(
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

func (fb *defaultFormBuilder) WriteField(fieldname, value string) error {
	return fb.writer.WriteField(fieldname, value)
}

func (fb *defaultFormBuilder) Close() error {
	return fb.writer.Close()
}

func (fb *defaultFormBuilder) FormDataContentType() string {
	return fb.writer.FormDataContentType()
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
