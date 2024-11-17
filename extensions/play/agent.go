package play

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	createAgentEndpoint   Endpoint = "/api/v1/agents"
	specificAgentEndpoint Endpoint = "/api/v1/agents/"
)

// AgentParams are the parameters for creating/updating an agent.
type AgentParams struct {
	Voice                           string `json:"voice"`
	VoiceSpeed                      string `json:"voiceSpeed"`
	DisplayName                     string `json:"displayName"`
	Description                     string `json:"description"`
	Greeting                        string `json:"greeting"`
	Prompt                          string `json:"prompt"`
	CriticalKnowledge               string `json:"criticalKnowledge"`
	Visibility                      string `json:"visibility"`
	AnswerOnlyFromCriticalKnowledge string `json:"answerOnlyFromCriticalKnowledge"`
	LLM                             string `json:"llm"`
}

// Agent represents an agent on the PlayAI platform.
type Agent struct {
	ID                string  `json:"id"`
	Voice             string  `json:"voice"`
	VoiceSpeed        float64 `json:"voiceSpeed"`
	DisplayName       string  `json:"displayName"`
	Description       string  `json:"description"`
	Greeting          string  `json:"greeting"`
	Prompt            string  `json:"prompt"`
	CriticalKnowledge string  `json:"criticalKnowledge"`
	Visibility        string  `json:"visibility"`
	Llm               struct {
		BaseURL    string `json:"baseURL"`
		APIKey     string `json:"apiKey"`
		BaseParams struct {
			DefaultHeaders struct {
				XMyExtraKey string `json:"x-my-extra-key"`
			} `json:"defaultHeaders"`
			Model       string `json:"model"`
			Temperature int    `json:"temperature"`
			MaxTokens   int    `json:"maxTokens"`
		} `json:"baseParams"`
	} `json:"llm"`
	AnswerOnlyFromCriticalKnowledge bool   `json:"answerOnlyFromCriticalKnowledge"`
	AvatarPhotoURL                  string `json:"avatarPhotoUrl"`
	CriticalKnowledgeFiles          []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
		Size int    `json:"size"`
		Type string `json:"type"`
	} `json:"criticalKnowledgeFiles"`
	PhoneNumbers []struct {
		PhoneNumber int64  `json:"phoneNumber"`
		Country     string `json:"country"`
		Locality    string `json:"locality"`
	} `json:"phoneNumbers"`
}

// CreateAgent creates a new agent.
func (p *PlayAI) CreateAgent(request AgentParams) (Agent, error) {
	req, err := builders.NewRequest(
		context.Background(),
		p.header,
		http.MethodPost,
		p.baseURL+string(createAgentEndpoint),
		builders.WithBody(request),
	)
	if err != nil {
		return Agent{}, err
	}
	var resp Agent
	err = p.sendRequest(req, &resp)
	if err != nil {
		return Agent{}, err
	}
	return resp, nil
}

// GetAgent gets an agent by id.
func (p *PlayAI) GetAgent(id string) (Agent, error) {
	req, err := builders.NewRequest(
		context.Background(),
		p.header,
		http.MethodGet,
		p.baseURL+string(specificAgentEndpoint)+id,
	)
	if err != nil {
		return Agent{}, err
	}
	var resp Agent
	err = p.sendRequest(req, &resp)
	if err != nil {
		return Agent{}, err
	}
	return resp, nil
}

// UpdateAgent updates an agent by id.
func (p *PlayAI) UpdateAgent(id string, request AgentParams) (Agent, error) {
	req, err := builders.NewRequest(
		context.Background(),
		p.header,
		http.MethodPut,
		p.baseURL+string(specificAgentEndpoint)+id,
		builders.WithBody(request),
	)
	if err != nil {
		return Agent{}, err
	}
	var resp Agent
	err = p.sendRequest(req, &resp)
	if err != nil {
		return Agent{}, err
	}
	return resp, nil
}

//TODO: Implemnt https://docs.play.ai/documentation/get-started/bring-your-own-llm
