// Package main is an example of using groq-go to create a chat completion
// using the llama-70B-tools-preview model to create headers for vhdl projects.
package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/groq-go"
)

var (
	//go:embed template.tmpl
	emTempl        string
	codeTemplate   *template.Template
	headerTemplate *template.Template
)

func init() {
	var err error
	codeTemplate, err = template.New("code").Parse(emTempl)
	if err != nil {
		log.Fatal(err)
	}
	headerTemplate, err = template.New("header").
		Funcs(template.FuncMap{}).
		Parse(emTempl)
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	ctx := context.Background()
	err := run(ctx, os.Getenv)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to run: %w", err))
	}
}
func run(
	ctx context.Context,
	getenv func(string) string,
) error {
	client, err := groq.NewClient(getenv("GROQ_KEY"))
	if err != nil {
		return err
	}
	log.Debugf("running with %s", getenv("GROQ_KEY"))
	for _, val := range fileMap {
	retry:
		time.Sleep(6 * time.Second)
		log.Debugf("processing %s", val.Destination)
		filename := strings.Split(val.Destination, "/")[len(strings.Split(val.Destination, "/"))-1]
		prompt, err := executeCodeTemplate(CodeTemplateData{
			Source: val.Source,
			Name:   filename,
			Files:  fileMap.ToFileArray(val.Destination),
		})
		if err != nil {
			return err
		}
		var thoughtThroughCode thoughtThroughCode
		err = client.CreateChatCompletionJSON(
			ctx,
			groq.ChatCompletionRequest{
				Model: groq.Llama3Groq70B8192ToolUsePreview,
				Messages: []groq.ChatCompletionMessage{{
					Role:    groq.ChatMessageRoleSystem,
					Content: prompt,
				}},
			},
			&thoughtThroughCode,
		)
		if err != nil {
			goto retry
		}
		log.Debugf(
			"thoughts for %s:  %s",
			val.Destination,
			thoughtThroughCode.Thoughts,
		)
		content, err := executeHeaderTemplate(HeaderTemplateData{
			FileName:    filename,
			Description: wrapText(thoughtThroughCode.Description),
			Code:        val.Source,
		})
		if err != nil {
			return err
		}
		log.Debugf("Creating file %s", val.Destination)
		oF, err := os.Create(val.Destination)
		if err != nil {
			return err
		}
		defer oF.Close()
		log.Debugf("Writing file %s", val.Destination)
		_, err = oF.WriteString(content)
		if err != nil {
			return err
		}
	}
	return nil
}

// FileMapper is the map of a source file content to a output/report folder.
type FileMapper []struct {
	Source      string
	Destination string
}

const (
	// Destinations
	muxDest             = "./report/Mux/"
	fullAdderDest       = "./report/Adder/"
	adderSubtractorDest = "./report/AddSub/"
	nMuxDest            = "./report/NMux/"
	onesCompDest        = "./report/OnesComp/"
	tpuElementDest      = "./report/MAC/"
)

var (
	fileMap = FileMapper{
		{Source: tbMux2t1, Destination: muxDest + "tb_Mux2t1.vhd"},
		{Source: mux2t1, Destination: muxDest + "mux2t1.vhd"},
		{Source: mux2t1s, Destination: muxDest + "mux2t1s.vhd"},
		{Source: tbMux2t1s, Destination: muxDest + "tb_Mux2t1s.vhd"},
		{Source: mux2t1N, Destination: nMuxDest + "mux2t1_N.vhd"},
		{Source: tbMux2t1, Destination: nMuxDest + "tb_Mux2t1.vhd"},
		{Source: mux2t1, Destination: nMuxDest + "mux2t1.vhd"},
		{Source: tbNMux2t1, Destination: nMuxDest + "tb_NMux2t1.vhd"},
		{Source: fullAdder, Destination: fullAdderDest + "FullAdder.vhd"},
		{Source: tbFullAdder, Destination: fullAdderDest + "tb_FullAdder.vhd"},
		{Source: tbNBitAdder, Destination: fullAdderDest + "tb_NBitAdder.vhd"},
		{Source: nBitAdder, Destination: fullAdderDest + "nBitAdder.vhd"},
		{
			Source:      adderSubtractor,
			Destination: adderSubtractorDest + "AdderSubtractor.vhd",
		},
		{
			Source:      tbAdderSubtractor,
			Destination: adderSubtractorDest + "tb_AdderSubtractor.vhd",
		},
		{
			Source:      nBitInverter,
			Destination: adderSubtractorDest + "nBitInverter.vhd",
		},
		{
			Source:      tbNBitInverter,
			Destination: adderSubtractorDest + "tb_nBitInverter.vhd",
		},
		{Source: mux2t1N, Destination: adderSubtractorDest + "mux2t1_N.vhd"},
		{
			Source:      nBitAdder,
			Destination: adderSubtractorDest + "nBitAdder.vhd",
		},
		{Source: xorg2, Destination: onesCompDest + "xorg2.vhd"},
		{Source: org2, Destination: onesCompDest + "org2.vhd"},
		{Source: onesComp, Destination: onesCompDest + "OnesComp.vhd"},
		{Source: tbOnesComp, Destination: onesCompDest + "tb_OnesComp.vhd"},
		{
			Source:      tpuMvElement,
			Destination: tpuElementDest + "tpuMvElement.vhd",
		},
		{Source: xorg2, Destination: tpuElementDest + "xorg2.vhd"},
		{Source: org2, Destination: tpuElementDest + "org2.vhd"},
		{
			Source:      nBitInverter,
			Destination: tpuElementDest + "nBitInverter.vhd",
		},
		{Source: nBitAdder, Destination: tpuElementDest + "nBitAdder.vhd"},
		{Source: mux2t1, Destination: tpuElementDest + "mux2t1.vhd"},
		{Source: andg2, Destination: tpuElementDest + "andg2.vhd"},
		{Source: regLd, Destination: tpuElementDest + "regLd.vhd"},
		{Source: invg, Destination: tpuElementDest + "invg.vhd"},
		{Source: adder, Destination: tpuElementDest + "adder.vhd"},
		{
			Source:      adderSubtractor,
			Destination: tpuElementDest + "adderSubtractor.vhd",
		},
		{Source: fullAdder, Destination: tpuElementDest + "fullAdder.vhd"},
		{Source: multiplier, Destination: tpuElementDest + "multiplier.vhd"},
		{Source: regLd, Destination: tpuElementDest + "regLd.vhd"},
		{Source: reg, Destination: tpuElementDest + "reg.vhd"},
		{
			Source:      tbTPUElement,
			Destination: tpuElementDest + "tb_TPUElement.vhd",
		},
	}
)

//go:embed src/Adder.vhd
var adder string

//go:embed src/AdderSubtractor.vhd
var adderSubtractor string

//go:embed src/FullAdder.vhd
var fullAdder string

//go:embed src/Multiplier.vhd
var multiplier string

//go:embed src/NBitAdder.vhd
var nBitAdder string

//go:embed src/NBitInverter.vhd
var nBitInverter string

//go:embed src/OnesComp.vhd
var onesComp string

//go:embed src/Reg.vhd
var reg string

//go:embed src/RegLd.vhd
var regLd string

//go:embed src/TPU_MV_Element.vhd
var tpuMvElement string

//go:embed src/andg2.vhd
var andg2 string

//go:embed src/invg.vhd
var invg string

//go:embed src/mux2t1.vhd
var mux2t1 string

//go:embed src/mux2t1s.vhd
var mux2t1s string

//go:embed src/mux2t1_N.vhd
var mux2t1N string

//go:embed src/org2.vhd
var org2 string

//go:embed src/xorg2.vhd
var xorg2 string

//go:embed test/tb_NMux2t1.vhd
var tbNMux2t1 string

//go:embed test/tb_NMux2t1.vhd
var tbMux2t1 string

//go:embed test/tb_nBitAdder.vhd
var tbNBitAdder string

//go:embed test/tb_mux2t1s.vhd
var tbMux2t1s string

//go:embed test/tb_nBitInverter.vhd
var tbNBitInverter string

//go:embed test/tb_AdderSubtractor.vhd
var tbAdderSubtractor string

//go:embed test/tb_NFullAdder.vhd
var tbFullAdder string

//go:embed test/tb_OnesComp.vhd
var tbOnesComp string

//go:embed test/tb_TPU_MV_Element.vhd
var tbTPUElement string

// File represents a single file with its name and content
type File struct {
	Name    string
	Content string
}

// ToFileArray converts the FileMapper to a File array
func (fM *FileMapper) ToFileArray(dest string) []File {
	var files []File
	for _, file := range *fM {
		split := strings.Split(file.Destination, "/")
		if strings.Contains(dest, "/"+split[2]+"/") {
			log.Debugf("adding %s to files", file.Destination)
			files = append(files, File{
				Name:    file.Destination,
				Content: file.Source,
			})
			continue
		}
	}
	log.Debugf("n-files: %v", len(files))
	return files
}

// CodeTemplateData represents the data structure for the code template
type CodeTemplateData struct {
	Source string
	Name   string
	Files  []File
}

// HeaderTemplateData represents the data structure for the header template
type HeaderTemplateData struct {
	FileName    string
	Description string
	Code        string
}
type thoughtThroughCode struct {
	Thoughts    string `json:"thoughts"    jsonschema:"title=Thoughts,description=Thoughts on the code and thinking through exactly how it interacts with other given code in the project."`
	Description string `json:"description" jsonschema:"title=Description,description=A description of the code's function, form, etc."`
}

func executeCodeTemplate(data CodeTemplateData) (string, error) {
	buf := new(bytes.Buffer)
	err := codeTemplate.Execute(buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute code template: %v", err)
	}
	return buf.String(), nil
}
func executeHeaderTemplate(data HeaderTemplateData) (string, error) {
	var result strings.Builder
	err := headerTemplate.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute header template: %v", err)
	}
	return result.String(), nil
}

// wrapText trims the string to 80 characters per line,
// adding newlines and hyphens if it is longer.
func wrapText(s string) string {
	var result strings.Builder
	maxLineLength := 80
	words := strings.Fields(s)
	lineLength := 0
	for i, word := range words {
		wordLength := len(word)
		if lineLength+wordLength > maxLineLength {
			if wordLength > maxLineLength {
				remaining := word
				for len(remaining) > 0 {
					spaceLeft := maxLineLength - lineLength
					if spaceLeft <= 0 {
						result.WriteString("\n-- ")
						lineLength = 0
						spaceLeft = maxLineLength
					}
					take := min(spaceLeft, len(remaining))
					if take < len(remaining) {
						// Add hyphen when breaking a word
						result.WriteString(remaining[:take] + "-\n")
						lineLength = 0
					} else {
						result.WriteString(remaining[:take])
						lineLength += take
					}
					remaining = remaining[take:]
				}
			} else {
				// Start a new line
				result.WriteString("\n-- ")
				result.WriteString(word)
				lineLength = wordLength
			}
		} else {
			if i > 0 {
				result.WriteString(" ")
				lineLength++
			}
			result.WriteString(word)
			lineLength += wordLength
		}
	}
	return result.String()
}
