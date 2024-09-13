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
	outputFile = "models.go"
)

//go:embed models.go.tmpl
var outputFileTemplate string

type templateParams struct {
	Models []ResponseModel `json:"models"`
}

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
	models := response.Data
	nameModels(models)
	err = fillTemplate(buf, models)
	if err != nil {
		return err
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf(
			"error formatting output: %w : %s",
			err,
			buf.String(),
		)
	}
	f, err := os.Create(outputFile)
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

func fillTemplate(w io.Writer, models []ResponseModel) error {
	funcMap := template.FuncMap{
		// isTextModel returns true if the model has a context window
		// greater than 1024 and is not specific models like
		// llama-guard-3-8b.
		"isTextModel": func(model ResponseModel) bool {
			if model.ID != "llama-guard-3-8b" {
				return model.ContextWindow > 1024
			}
			return false
		},
		// isAudioModel returns true if the model has a context window
		// less than 1024 and is not specific models like
		// llama-guard-3-8b.
		"isAudioModel": func(model ResponseModel) bool {
			if model.ID != "llama-guard-3-8b" {
				return model.ContextWindow < 1024
			}
			return false
		},
		// isModerationModel returns true if the model is
		// a model that can be used for moderation.
		//
		// llama-guard-3-8b is a moderation model.
		"isModerationModel": func(model ResponseModel) bool {
			// if the id of the model is llama-guard-3-8b
			return model.ID == "llama-guard-3-8b"
		},
		// getCurrentDate returns the current date in the format
		// "2006-01-02 15:04:05".
		"getCurrentDate": func() string {
			return time.Now().Format("2006-01-02 15:04:05")
		},
	}
	tmpla, err := template.New("models").
		Funcs(funcMap).
		Parse(outputFileTemplate)
	if err != nil {
		return err
	}
	err = tmpla.Execute(w, templateParams{Models: models})
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
