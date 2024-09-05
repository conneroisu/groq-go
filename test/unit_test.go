package test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"

	"github.com/conneroisu/groq-go"
	"github.com/stretchr/testify/assert"
)

func TestTestServer(t *testing.T) {
	num := rand.Intn(100)
	a := assert.New(t)
	ctx := context.Background()
	client, err := groq.NewClient(os.Getenv("GROQ_KEY"))
	a.NoError(err, "NewClient error")
	strm, err := client.CreateChatCompletionStream(
		ctx,
		groq.ChatCompletionRequest{
			Model: groq.Llama3070B8192ToolUsePreview,
			Messages: []groq.ChatCompletionMessage{
				{
					Role: groq.ChatMessageRoleUser,
					// convert the content of a excel file into a csv that can be imported into google calendar:
					// <?xml version="1.0" encoding="UTF-8"?>
					// <worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><dimension ref="A1:F7"/><sheetViews><sheetView workbookViewId="0" tabSelected="true" rightToLeft="false"/></sheetViews><sheetFormatPr defaultRowHeight="15.0" baseColWidth="8"/><cols><col min="1" max="1" width="23.0" customWidth="true"/><col min="2" max="2" width="23.0" customWidth="true"/><col min="3" max="3" width="23.0" customWidth="true"/><col min="4" max="4" width="23.0" customWidth="true"/><col min="5" max="5" width="23.0" customWidth="true"/><col min="6" max="6" width="23.0" customWidth="true"/></cols><sheetData>
					// <row r="1">
					// <c r="A1" s="1" t="inlineStr"><is><t>Course Listing</t></is></c><c r="B1" s="1" t="inlineStr"><is><t>Section</t></is></c><c r="C1" s="1" t="inlineStr"><is><t>Instructional Format</t></is></c><c r="D1" s="1" t="inlineStr"><is><t>Delivery Mode</t></is></c><c r="E1" s="1" t="inlineStr"><is><t>Meeting Patterns</t></is></c><c r="F1" s="1" t="inlineStr"><is><t>Instructor</t></is></c></row>
					// <row r="2">
					// <c r="A2" s="2" t="inlineStr"><is><t>CPRE 3810 - Computer Organization and Assembly Level Programming</t></is></c><c r="B2" s="2" t="inlineStr"><is><t>CPRE 3810-1 - Computer Organization and Assembly Level Programming</t></is></c><c r="C2" s="2" t="inlineStr"><is><t>Lecture</t></is></c><c r="D2" s="2" t="inlineStr"><is><t>In-Person</t></is></c><c r="E2" s="2" t="inlineStr"><is><t>MWF | 8:50 AM - 9:40 AM | 0101 CARVER - Carver Hall</t></is></c><c r="F2" s="2" t="inlineStr"><is><t>Berk Gulmezoglu</t></is></c></row>
					// <row r="3">
					// <c r="A3" s="2" t="inlineStr"><is><t>EE 3240 - Signals and Systems II</t></is></c><c r="B3" s="2" t="inlineStr"><is><t>EE 3240-1 - Signals and Systems II</t></is></c><c r="C3" s="2" t="inlineStr"><is><t>Lecture</t></is></c><c r="D3" s="2" t="inlineStr"><is><t>In-Person</t></is></c><c r="E3" s="2" t="inlineStr"><is><t>MWF | 2:15 PM - 3:05 PM | 1134 SWEENEY - Sweeney Hall</t></is></c><c r="F3" s="2" t="inlineStr"><is><t>Ratnesh Kumar</t></is></c></row>
					// <row r="4">
					// <c r="A4" s="2" t="inlineStr"><is><t>EE 3220 - Probabilistic Methods for Electrical Engineers</t></is></c><c r="B4" s="2" t="inlineStr"><is><t>EE 3220-1 - Probabilistic Methods for Electrical Engineers</t></is></c><c r="C4" s="2" t="inlineStr"><is><t>Lecture</t></is></c><c r="D4" s="2" t="inlineStr"><is><t>In-Person</t></is></c><c r="E4" s="2" t="inlineStr"><is><t>TR | 11:00 AM - 12:15 PM | 1134 SWEENEY - Sweeney Hall</t></is></c><c r="F4" s="2" t="inlineStr"><is><t>Julie A Dickerson</t></is></c></row>
					// <row r="5">
					// <c r="A5" s="2" t="inlineStr"><is><t>SOC 1340 - Introduction to Sociology</t></is></c><c r="B5" s="2" t="inlineStr"><is><t>SOC 1340-2 - Introduction to Sociology</t></is></c><c r="C5" s="2" t="inlineStr"><is><t>Lecture</t></is></c><c r="D5" s="2" t="inlineStr"><is><t>In-Person</t></is></c><c r="E5" s="2" t="inlineStr"><is><t>MWF | 11:00 AM - 11:50 AM | 0127 CURTISS - Curtiss Hall</t></is></c><c r="F5" s="2" t="inlineStr"><is><t>David Scott Schweingruber</t></is></c></row>
					// <row r="6">
					// <c r="A6" s="2" t="inlineStr"><is><t>CPRE 3810 - Computer Organization and Assembly Level Programming</t></is></c><c r="B6" s="2" t="inlineStr"><is><t>CPRE 3810-F - Computer Organization and Assembly Level Programming</t></is></c><c r="C6" s="2" t="inlineStr"><is><t>Laboratory</t></is></c><c r="D6" s="2" t="inlineStr"><is><t>In-Person</t></is></c><c r="E6" s="2" t="inlineStr"><is><t>R | 8:00 AM - 9:50 AM | 2050 COOVER - Coover Hall</t></is></c><c r="F6" s="2" t="inlineStr"><is><t>Berk Gulmezoglu</t></is></c></row>
					// <row r="7">
					// <c r="A7" s="2" t="inlineStr"><is><t>EE 3240 - Signals and Systems II</t></is></c><c r="B7" s="2" t="inlineStr"><is><t>EE 3240-C - Signals and Systems II</t></is></c><c r="C7" s="2" t="inlineStr"><is><t>Laboratory</t></is></c><c r="D7" s="2" t="inlineStr"><is><t>In-Person</t></is></c><c r="E7" s="2" t="inlineStr"><is><t>R | 4:10 PM - 7:00 PM | 2011 COOVER - Coover Hall</t></is></c><c r="F7" s="2" t="inlineStr"><is><t>Ratnesh Kumar</t></is></c></row>
					// </sheetData><pageMargins bottom="0.75" footer="0.3" header="0.3" left="0.7" right="0.7" top="0.75"/></worksheet>
					Content: fmt.Sprintf(`
problem: %d
You have a six-sided die that you roll once. Let $R{i}$ denote the event that the roll is $i$. Let $G{j}$ denote the event that the roll is greater than $j$. Let $E$ denote the event that the roll of the die is even-numbered.
(a) What is $P\left[R{3} \mid G{1}\right]$, the conditional probability that 3 is rolled given that the roll is greater than 1 ?
(b) What is the conditional probability that 6 is rolled given that the roll is greater than 3 ?
(c) What is the $P\left[G_{3} \mid E\right]$, the conditional probability that the roll is greater than 3 given that the roll is even?
(d) Given that the roll is greater than 3, what is the conditional probability that the roll is even?
					`, num,
					),
				},
			},
			MaxTokens: 2000,
			Stream:    true,
		},
	)
	a.NoError(err, "CreateCompletionStream error")

	i := 0
	for {
		i++
		val, err := strm.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		// t.Logf("%d %s\n", i, val.Choices[0].Delta.Content)
		print(val.Choices[0].Delta.Content)
	}
}
