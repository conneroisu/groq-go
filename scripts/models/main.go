package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	outputFile = "models.go"
)

// main is the entry point for the application.
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// Response is a response from the models endpoint.
type Response struct {
	Object string `json:"object"`
	Data   []struct {
		ID            string `json:"id"`
		Object        string `json:"object"`
		Created       int    `json:"created"`
		OwnedBy       string `json:"owned_by"`
		Active        bool   `json:"active"`
		ContextWindow int    `json:"context_window"`
		PublicApps    any    `json:"public_apps"`
	} `json:"data"`
}

// run runs the main function.
func run() error {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://api.groq.com/openai/v1/models",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY"))
	req.Header.Set("Content-Type", "application/json")
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
	var response Response
	err = json.Unmarshal(bodyText, &response)
	if err != nil {
		return err
	}
	return nil
}
