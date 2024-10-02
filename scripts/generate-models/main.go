// Package main is the main package for the groq-modeler.
//
// It is used to generate the models for the groq-go library.
package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"text/template"
	"time"

	"github.com/samber/lo"
)

const (
	modelFileName     = "models.go"
	modelTestFileName = "models_test.go"
)

var (
	modelTemplate = template.New("models").Funcs(funcMap)
	testTemplate  = template.New("test").Funcs(funcMap)

	//go:embed models.go.tmpl
	modelFileTemplate string
	//go:embed models_test.go.tmpl
	testFileTemplate string

	funcMap = template.FuncMap{
		"getCurrentDate": func() string {
			return time.Now().Format("2006-01-02 15:04:05")
		},
	}
)

// main is the entry point for the application.
func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

// Response is a response from the models endpoint.
type Response struct {
	Object string          `json:"object"`
	Data   []ResponseModel `json:"data"`
}

// ResponseModel is a response from the models endpoint.
type ResponseModel struct {
	ID            string `json:"id"`
	Name          string `json:"name,omitempty"`
	Object        string `json:"object"`
	Created       int    `json:"created"`
	OwnedBy       string `json:"owned_by"`
	Active        bool   `json:"active"`
	ContextWindow int    `json:"context_window"`
	PublicApps    any    `json:"public_apps"`
}

// CategorizedModels is a struct that contains all the models.
type CategorizedModels struct {
	ChatModels          []ResponseModel `json:"text"`
	TranscriptionModels []ResponseModel `json:"transcription"`
	TranslationModels   []ResponseModel `json:"translation"`
	ModerationModels    []ResponseModel `json:"moderation"`
	MultiModalModels    []ResponseModel `json:"multi_modal"`
}

// Categorize returns a Categorize struct with all the models.
func (r *Response) Categorize() (CategorizedModels, error) {
	var models CategorizedModels

	nameModels(r.Data)
	for _, model := range r.Data {
		if isTextModel(model) {
			models.ChatModels = append(models.ChatModels, model)
		}
		if isTranscriptionModel(model) {
			models.TranscriptionModels = append(models.TranscriptionModels, model)
		}
		if isTranslationModel(model) {
			models.TranslationModels = append(models.TranslationModels, model)
		}
		if isModerationModel(model) {
			models.ModerationModels = append(models.ModerationModels, model)
		}
		if isMultiModalModel(model) {
			models.MultiModalModels = append(models.MultiModalModels, model)
		}
	}
	return models, nil
}

func isMultiModalModel(model ResponseModel) bool {
	return false
}

func isTextModel(model ResponseModel) bool {
	if model.ID != "llama-guard-3-8b" {
		return model.ContextWindow > 1024
	}
	return false
}

func isModerationModel(model ResponseModel) bool {
	// if the id of the model is llama-guard-3-8b
	return model.ID == "llama-guard-3-8b"
}

func isTranslationModel(model ResponseModel) bool {
	return model.ID == "whisper-large-v3"
}

func isTranscriptionModel(model ResponseModel) bool {
	return model.ID == "whisper-large-v3"
}

// run runs the main function.
func run(_ context.Context) error {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://api.groq.com/openai/v1/models",
		nil,
	)
	if err != nil {
		return err
	}
	key := os.Getenv("GROQ_KEY")
	if key == "" {
		return fmt.Errorf("GROQ_KEY is not set")
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var response Response
	err = json.Unmarshal(bodyText, &response)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	ms, err := response.Categorize()
	if err != nil {
		return err
	}
	err = fillModelsTemplate(buf, ms)
	if err != nil {
		return err
	}
	formatted, err := cleanFile(buf)
	if err != nil {
		return err
	}
	f, err := os.Create(modelFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(formatted)
	if err != nil {
		return err
	}
	buf.Reset()
	err = fillTestTemplate(buf, ms)
	if err != nil {
		return err
	}
	formatted, err = cleanFile(buf)
	if err != nil {
		return err
	}
	f, err = os.Create(modelTestFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(formatted)
	if err != nil {
		return err
	}
	return nil
}

func cleanFile(r io.Reader) ([]byte, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	formatted, err := format.Source(b)
	if err != nil {
		return nil, fmt.Errorf(
			"error formatting output: %w : %s",
			err,
			b,
		)
	}
	return formatted, nil
}

func fillModelsTemplate(w io.Writer, models CategorizedModels) (err error) {
	modelTemplate, err = modelTemplate.Parse(modelFileTemplate)
	if err != nil {
		return err
	}
	err = modelTemplate.Execute(w, models)
	if err != nil {
		return err
	}
	return nil
}

func nameModels(models []ResponseModel) {
	for i := range models {
		if (models)[i].Name == "" {
			models[i].Name = lo.PascalCase(models[i].ID)
		}
	}
	// sort models by name alphabetically
	sort.Slice(models, func(i, j int) bool {
		return models[i].Name < models[j].Name
	})
}

func fillTestTemplate(w io.Writer, models CategorizedModels) (err error) {
	testTemplate, err = testTemplate.Parse(testFileTemplate)
	if err != nil {
		return err
	}
	err = testTemplate.Execute(w, models)
	if err != nil {
		return err
	}
	return nil
}
