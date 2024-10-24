package composio

// https://backend.composio.dev/api/v1/connectedAccounts?user_uuid=default&showActiveOnly=true

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/conneroisu/groq-go/pkg/builders"
)

type (
	// Auther is an interface for composio auth.
	Auther interface {
		GetConnectedAccounts(ctx context.Context, opts ...AuthOption) (ConnectedAccounts, error)
	}
	// ConnectedAccounts represents a composio connected account.
	ConnectedAccounts struct {
		Items []struct {
			IntegrationID    string `json:"integrationId"`
			ConnectionParams struct {
				Scope            string   `json:"scope"`
				Scopes           []string `json:"scopes"`
				BaseURL          string   `json:"base_url"`
				ClientID         string   `json:"client_id"`
				TokenType        string   `json:"token_type"`
				RedirectURL      string   `json:"redirectUrl"`
				AccessToken      string   `json:"access_token"`
				CallbackURL      string   `json:"callback_url"`
				ClientSecret     string   `json:"client_secret"`
				CodeVerifier     string   `json:"code_verifier"`
				FinalRedirectURI string   `json:"finalRedirectUri"`
			} `json:"connectionParams"`
			IsDisabled         bool      `json:"isDisabled"`
			ID                 string    `json:"id"`
			MemberID           string    `json:"memberId"`
			ClientUniqueUserID string    `json:"clientUniqueUserId"`
			Status             string    `json:"status"`
			Enabled            bool      `json:"enabled"`
			CreatedAt          time.Time `json:"createdAt"`
			UpdatedAt          time.Time `json:"updatedAt"`
			Member             struct {
				ID        string    `json:"id"`
				ClientID  string    `json:"clientId"`
				Email     string    `json:"email"`
				Name      string    `json:"name"`
				Role      string    `json:"role"`
				Metadata  any       `json:"metadata"`
				CreatedAt time.Time `json:"createdAt"`
				UpdatedAt time.Time `json:"updatedAt"`
				DeletedAt any       `json:"deletedAt"`
			} `json:"member"`
			AppUniqueID               string `json:"appUniqueId"`
			AppName                   string `json:"appName"`
			Logo                      string `json:"logo"`
			IntegrationIsDisabled     bool   `json:"integrationIsDisabled"`
			IntegrationDisabledReason string `json:"integrationDisabledReason"`
			InvocationCount           string `json:"invocationCount"`
		} `json:"items"`
		TotalPages int `json:"totalPages"`
		Page       int `json:"page"`
	}
)

// GetConnectedAccounts returns the connected accounts for the composio client.
func (c *Composio) GetConnectedAccounts(ctx context.Context, opts ...AuthOption) (ca ConnectedAccounts, err error) {
	uri := fmt.Sprintf("%s/v1/connectedAccounts", c.baseURL)
	u, err := url.Parse(uri)
	if err != nil {
		return ca, err
	}
	ps := url.Values{}
	ps.Add("user_uuid", "default")
	ps.Add("showActiveOnly", "true")
	for _, opt := range opts {
		opt(u)
	}
	u.RawQuery = ps.Encode()
	uri = u.String()
	c.logger.Debug("auth", "url", uri)
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodGet,
		uri,
		builders.WithBody(nil),
	)
	if err != nil {
		return ca, err
	}
	err = c.doRequest(req, &ca)
	return ca, err
}
