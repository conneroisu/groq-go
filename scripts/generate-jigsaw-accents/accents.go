package main

import (
	"bytes"
	_ "embed"
	"text/template"

	_ "embed"
)

//go:embed accents.go.tmpl
var accentsTemplate string

var (
	textTemplate = template.Must(template.New("accents").Parse(accentsTemplate))
)

func FillAccents(accents SpeakerVoiceAccent) string {
	var buf bytes.Buffer
	err := textTemplate.Execute(&buf, accents)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
