//go:build !test
// +build !test

package groq

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestAudioWithFailingFormBuilder(t *testing.T) {
	a := assert.New(t)
	dir, cleanup := test.CreateTestDirectory(t)
	defer cleanup()
	path := filepath.Join(dir, "fake.mp3")
	test.CreateTestFile(t, path)

	req := AudioRequest{
		FilePath:    path,
		Prompt:      "test",
		Temperature: 0.5,
		Language:    "en",
		Format:      AudioResponseFormatSRT,
	}

	mockFailedErr := fmt.Errorf("mock form builder fail")
	mockBuilder := &mockFormBuilder{}

	mockBuilder.mockCreateFormFile = func(string, *os.File) error {
		return mockFailedErr
	}
	err := audioMultipartForm(req, mockBuilder)
	a.ErrorIs(
		err,
		mockFailedErr,
		"audioMultipartForm should return error if form builder fails",
	)

	mockBuilder.mockCreateFormFile = func(string, *os.File) error {
		return nil
	}

	var failForField string
	mockBuilder.mockWriteField = func(fieldname, _ string) error {
		if fieldname == failForField {
			return mockFailedErr
		}
		return nil
	}

	failOn := []string{
		"model",
		"prompt",
		"temperature",
		"language",
		"response_format",
	}
	for _, failingField := range failOn {
		failForField = failingField
		mockFailedErr = fmt.Errorf(
			"mock form builder fail on field %s",
			failingField,
		)

		err = audioMultipartForm(req, mockBuilder)
		a.Error(
			err,
			mockFailedErr,
			"audioMultipartForm should return error if form builder fails",
		)
	}
}

func TestCreateFileField(t *testing.T) {
	a := assert.New(t)
	t.Run("createFileField failing file", func(t *testing.T) {
		t.Parallel()
		dir, cleanup := test.CreateTestDirectory(t)
		defer cleanup()
		path := filepath.Join(dir, "fake.mp3")
		test.CreateTestFile(t, path)
		req := AudioRequest{
			FilePath: path,
		}
		mockFailedErr := fmt.Errorf("mock form builder fail")
		mockBuilder := &mockFormBuilder{
			mockCreateFormFile: func(string, *os.File) error {
				return mockFailedErr
			},
		}
		err := createFileField(req, mockBuilder)
		a.ErrorIs(
			err,
			mockFailedErr,
			"createFileField using a file should return error if form builder fails",
		)
	})

	t.Run("createFileField failing reader", func(t *testing.T) {
		t.Parallel()
		req := AudioRequest{
			FilePath: "test.wav",
			Reader:   bytes.NewBuffer([]byte(`wav test contents`)),
		}

		mockFailedErr := fmt.Errorf("mock form builder fail")
		mockBuilder := &mockFormBuilder{
			mockCreateFormFileReader: func(string, io.Reader, string) error {
				return mockFailedErr
			},
		}

		err := createFileField(req, mockBuilder)
		a.ErrorIs(
			err,
			mockFailedErr,
			"createFileField using a reader should return error if form builder fails",
		)
	})

	t.Run("createFileField failing open", func(t *testing.T) {
		t.Parallel()
		req := AudioRequest{
			FilePath: "non_existing_file.wav",
		}
		mockBuilder := builders.NewFormBuilder(&test.FailingErrorBuffer{})
		err := createFileField(req, mockBuilder)
		a.Error(
			err,
			"createFileField using file should return error when open file fails",
		)
	})
}

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
