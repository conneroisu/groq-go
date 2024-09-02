package groq //nolint:testpackage // testing private field

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/conneroisu/groq-go/internal/test"
)

type mockFormBuilder struct {
	mockCreateFormFile       func(string, *os.File) error
	mockCreateFormFileReader func(string, io.Reader, string) error
	mockWriteField           func(string, string) error
	mockClose                func() error
}

func (fb *mockFormBuilder) CreateFormFile(fieldname string, file *os.File) error {
	return fb.mockCreateFormFile(fieldname, file)
}

func (fb *mockFormBuilder) CreateFormFileReader(fieldname string, r io.Reader, filename string) error {
	return fb.mockCreateFormFileReader(fieldname, r, filename)
}

func (fb *mockFormBuilder) WriteField(fieldname, value string) error {
	return fb.mockWriteField(fieldname, value)
}

func (fb *mockFormBuilder) Close() error {
	return fb.mockClose()
}

func (fb *mockFormBuilder) FormDataContentType() string {
	return ""
}

type failingWriter struct {
}

var errMockFailingWriterError = errors.New("mock writer failed")

func (*failingWriter) Write([]byte) (int, error) {
	return 0, errMockFailingWriterError
}

func TestFormBuilderWithFailingWriter(t *testing.T) {
	dir, cleanup := test.CreateTestDirectory(t)
	defer cleanup()

	file, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Errorf("Error creating tmp file: %v", err)
	}
	defer file.Close()
	defer os.Remove(file.Name())

	builder := NewFormBuilder(&failingWriter{})
	err = builder.CreateFormFile("file", file)
	a.ErrorIs(t, err, errMockFailingWriterError, "formbuilder should return error if writer fails")
}

func TestFormBuilderWithClosedFile(t *testing.T) {
	dir, cleanup := test.CreateTestDirectory(t)
	defer cleanup()

	file, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Errorf("Error creating tmp file: %v", err)
	}
	file.Close()
	defer os.Remove(file.Name())

	body := &bytes.Buffer{}
	builder := NewFormBuilder(body)
	err = builder.CreateFormFile("file", file)
	a.HasError(t, err, "formbuilder should return error if file is closed")
	a.ErrorIs(t, err, os.ErrClosed, "formbuilder should return error if file is closed")
}
