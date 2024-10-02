package e2b

import (
	"github.com/conneroisu/groq-go/pkg/mime"
)

type (
	APIError struct {
		Name      string
		Value     string
		Traceback []string
	}
	KernelException struct {
		Message string
	}
	Client struct {
		apiKey string
	}
	JupyterResult struct {
		Text         string            `json:"text"`
		HTML         string            `json:"html"`
		Markdown     string            `json:"markdown"`
		SVG          string            `json:"svg"`
		PNG          string            `json:"png"`
		JPEG         string            `json:"jpeg"`
		PDF          string            `json:"pdf"`
		LaTeX        string            `json:"latex"`
		JSON         string            `json:"json"`
		JavaScript   string            `json:"javascript"`
		Extra        map[string]string `json:"extra"`
		IsMainResult bool              `json:"is_main_result"`
		raw          map[mime.Type]string
	}

	Logs struct {
		Stdout []string `json:"stdout"`
		Stderr []string `json:"stderr"`
	}
	Execution struct {
		Results        []JupyterResult `json:"results"`
		Logs           Logs            `json:"logs"`
		Error          *APIError       `json:"error"`
		ExecutionCount int             `json:"execution_count"`
	}
)
