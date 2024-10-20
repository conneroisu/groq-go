package composio

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/conneroisu/groq-go"
)

const (
	composioAPIURLv1 = "https://api.composio.com/v1"
)

type (
	Composio struct {
		apiKey string
	}
	Composer interface {
		GetTools() []groq.Tool
	}
)

func (c *Composio) GetTools() []groq.Tool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://backend.composio.dev/api/v2/actions", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-API-Key", c.apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
	var tools []groq.Tool
	err = json.Unmarshal(bodyText, &tools)
	if err != nil {
		log.Fatal(err)
	}
	return tools
}
