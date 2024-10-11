package builders

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
)

type (
	// FormBuilder is an interface for building forms.
	FormBuilder interface {
		io.Closer
		CreateFormFile(fieldname string, file *os.File) error
		CreateFormFileReader(fieldname string, r io.Reader, filename string) error
		WriteField(fieldname, value string) error
		FormDataContentType() string
	}
	defaultFormBuilder struct {
		writer *multipart.Writer
	}
)

// NewFormBuilder creates a new DefaultFormBuilder.
func NewFormBuilder(body io.Writer) FormBuilder {
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
