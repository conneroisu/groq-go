// Package main is a script to generate the jigsaw accents for the project.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	_ "embed"

	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/conneroisu/groq-go/pkg/test"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	defaultBaseURL  = "https://api.jigsawstack.com"
	accentsEndpoint = "/v1/ai/tts"
)

type (
	// Client is a JigsawStack extension.
	Client struct {
		baseURL string
		apiKey  string
		client  *http.Client
		logger  *slog.Logger
		header  builders.Header
	}
	// SpeakerVoiceAccent represents a speaker voice accent.
	SpeakerVoiceAccent struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Accents []struct {
			GoName     string `json:"go_name"`
			Accent     string `json:"accent"`
			LocaleName string `json:"locale_name"`
			Gender     string `json:"gender"`
		} `json:"accents"`
	}
)

// AudioGetSpeakerVoiceAccents gets the speaker voice accents.
//
// GET https://api.jigsawstack.com/v1/ai/tts
//
// https://docs.jigsawstack.com/api-reference/audio/speaker-voice-accents
func (c *Client) AudioGetSpeakerVoiceAccents(
	ctx context.Context,
) (response SpeakerVoiceAccent, err error) {
	uri := c.baseURL + accentsEndpoint
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodGet,
		uri,
	)
	if err != nil {
		return
	}
	var resp SpeakerVoiceAccent
	err = c.sendRequest(req, &resp)
	if err != nil {
		return
	}
	if !resp.Success {
		return resp, fmt.Errorf("failed to get accents: %v", resp.Message)
	}
	for i := range resp.Accents {
		accent := &resp.Accents[i]
		accent.GoName = PascalCase(strings.ReplaceAll(accent.Accent, "-", ""))
	}
	return resp, nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Println(err)
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func newClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		header: builders.Header{SetCommonHeaders: func(r *http.Request) {
			r.Header.Set("x-api-key", apiKey)
		}},
		client:  http.DefaultClient,
		baseURL: defaultBaseURL,
	}
}

func run(ctx context.Context) error {
	println("generating accents")
	key, err := test.GetAPIKey("JIGSAWSTACK_API_KEY")
	if err != nil {
		return err
	}
	client := newClient(key)
	accents, err := client.AudioGetSpeakerVoiceAccents(ctx)
	if err != nil {
		return err
	}
	output := FillAccents(accents)
	println(output)
	return nil
}

func (c *Client) sendRequest(req *http.Request, v any) error {
	req.Header.Set("Accept", "application/json")
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK ||
		res.StatusCode >= http.StatusBadRequest {
		return nil
	}
	if v == nil {
		return nil
	}
	switch o := v.(type) {
	case *string:
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		*o = string(b)
		return nil
	default:
		err = json.NewDecoder(res.Body).Decode(v)
		if err != nil {
			read, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			c.logger.Debug("failed to decode response", "response", string(read))
			return fmt.Errorf("failed to decode response: %w\nbody: %s", err, string(read))
		}
		return nil
	}
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

//go:embed accents.go.tmpl
var accentsTemplate string

var (
	textTemplate = template.Must(template.New("accents").Parse(accentsTemplate))
)

// FillAccents fills the accents template with the given accents
func FillAccents(accents SpeakerVoiceAccent) string {
	var buf bytes.Buffer
	err := textTemplate.Execute(&buf, accents)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
