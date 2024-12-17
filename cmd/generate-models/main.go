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
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type (
	// Response is a response from the models endpoint.
	Response struct {
		Object string          `json:"object"`
		Data   []ResponseModel `json:"data"`
	}
	// ResponseModel is a response from the models endpoint.
	ResponseModel struct {
		ID            string `json:"id"`
		Name          string `json:"name,omitempty"`
		Object        string `json:"object"`
		Created       int    `json:"created"`
		OwnedBy       string `json:"owned_by"`
		Active        bool   `json:"active"`
		ContextWindow int    `json:"context_window"`
	}
	// CategorizedModels is a struct that contains all the models.
	CategorizedModels struct {
		ChatModels       []ResponseModel
		AudioModels      []ResponseModel
		ModerationModels []ResponseModel
	}
)

var (
	//go:embed template.tmpl
	templateContent string

	outputs = []*template.Template{
		template.New("models").Funcs(funcMap),
		template.New("models_test").Funcs(funcMap),
	}

	funcMap = template.FuncMap{
		"getCurrentDate": func() string {
			return time.Now().Format("2006-01-02 15:04:05")
		}}
)

// main is the entry point for the application.
func main() {
	if err := run(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// run runs the main function.
func run(ctx context.Context) error {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(
		ctx,
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
	ms, err := response.Categorize()
	if err != nil {
		return err
	}
	for _, template := range outputs {
		buf := new(bytes.Buffer)
		template, err = template.Parse(templateContent)
		if err != nil {
			return err
		}
		err = template.Execute(buf, ms)
		if err != nil {
			return err
		}
		messy, err := io.ReadAll(buf)
		if err != nil {
			return err
		}
		formatted, err := format.Source(messy)
		if err != nil {
			return err
		}
		f, err := os.Create("./" + template.Name() + ".go")
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(formatted)
		if err != nil {
			return err
		}
	}
	return nil
}

// Categorize returns a Categorize struct with all the models.
func (r *Response) Categorize() (CategorizedModels, error) {
	var models CategorizedModels
	for i := range r.Data {
		if (r.Data)[i].Name == "" {
			r.Data[i].Name = PascalCase(r.Data[i].ID)
		}
	}
	// sort models by name alphabetically
	sort.Slice(r.Data, func(i, j int) bool {
		return r.Data[i].Name < r.Data[j].Name
	})
	for _, model := range r.Data {
		if model.ID == "llama-guard-3-8b" {
			models.ModerationModels = append(models.ModerationModels, model)
			continue
		}
		if model.ContextWindow >= 1024 {
			models.ChatModels = append(models.ChatModels, model)
			continue
		}
		models.AudioModels = append(models.AudioModels, model)
	}
	return models, nil
}

var (
	// LowerCaseLettersCharset is a set of lower case letters.
	LowerCaseLettersCharset = []rune("abcdefghijklmnopqrstuvwxyz")
	// UpperCaseLettersCharset is a set of upper case letters.
	UpperCaseLettersCharset = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// LettersCharset is a set of letters.
	LettersCharset = append(LowerCaseLettersCharset, UpperCaseLettersCharset...)
	// NumbersCharset is a set of numbers.
	NumbersCharset = []rune("0123456789")
	// AlphanumericCharset is a set of alphanumeric characters.
	AlphanumericCharset = append(LettersCharset, NumbersCharset...)
	// SpecialCharset is a set of special characters.
	SpecialCharset = []rune("!@#$%^&*()_+-=[]{}|;':\",./<>?")
	// AllCharset is a set of all characters.
	AllCharset = append(AlphanumericCharset, SpecialCharset...)

	// bearer:disable go_lang_permissive_regex_validation
	splitWordReg = regexp.MustCompile(`([a-z])([A-Z0-9])|([a-zA-Z])([0-9])|([0-9])([a-zA-Z])|([A-Z])([A-Z])([a-z])`)
	// bearer:disable go_lang_permissive_regex_validation
	splitNumberLetterReg = regexp.MustCompile(`([0-9])([a-zA-Z])`)
)

// Words splits string into an array of its words.
func Words(str string) []string {
	str = splitWordReg.ReplaceAllString(str, `$1$3$5$7 $2$4$6$8$9`)
	// example: Int8Value => Int 8Value => Int 8 Value
	str = splitNumberLetterReg.ReplaceAllString(str, "$1 $2")
	var result strings.Builder
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune(' ')
		}
	}
	return strings.Fields(result.String())
}

// Capitalize converts the first character of string to upper case and the remaining to lower case.
func Capitalize(str string) string {
	return cases.Title(language.English).String(str)
}

// PascalCase converts string to pascal case.
func PascalCase(str string) string {
	items := Words(str)
	for i := range items {
		items[i] = Capitalize(items[i])
	}
	return strings.Join(items, "")
}
