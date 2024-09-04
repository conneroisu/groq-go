package groq //nolint:testpackage // testing private field

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/conneroisu/groq-go/internal/test"
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
		TimestampGranularities: []TranscriptionTimestampGranularity{
			TranscriptionTimestampGranularitySegment,
			TranscriptionTimestampGranularityWord,
		},
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
		"timestamp_granularities[]",
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
		mockBuilder := &mockFormBuilder{}
		err := createFileField(req, mockBuilder)
		a.Error(
			err,
			"createFileField using file should return error when open file fails",
		)
	})
}
