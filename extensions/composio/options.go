package composio

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

type (
	// Option is an option for the composio client.
	//
	// WithLogger sets the logger for the composio client.
	Option func(*Composio)

	// ToolsOption is an option for the tools request.
	ToolsOption func(*url.Values)

	// AuthOption is an option for the auth request.
	AuthOption func(*url.Values)
)

// Composer Options

// WithLogger sets the logger for the composio client.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Composio) { c.logger = logger }
}

// WithBaseURL sets the base URL for the composio client.
func WithBaseURL(baseURL string) Option {
	return func(c *Composio) { c.baseURL = baseURL }
}

// Get Tool Options

// WithTags sets the tags for the tools request.
func WithTags(tags ...string) ToolsOption {
	return func(u *url.Values) { u.Add("tags", strings.Join(tags, ",")) }
}

// WithApp sets the app for the tools request.
func WithApp(app string) ToolsOption {
	return func(u *url.Values) { u.Add("appNames", app) }
}

// WithEntityID sets the entity id for the tools request.
func WithEntityID(entityID string) ToolsOption {
	return func(u *url.Values) { u.Add("user_uuid", entityID) }
}

// WithUseCase sets the use case for the tools request.
func WithUseCase(useCase string) ToolsOption {
	return func(u *url.Values) { u.Add("useCase", useCase) }
}

// Auth Options

// WithShowActiveOnly sets the show active only for the auth request.
func WithShowActiveOnly(showActiveOnly bool) AuthOption {
	return func(u *url.Values) {
		u.Set("showActiveOnly", fmt.Sprintf("%t", showActiveOnly))
	}
}

// WithUserUUID sets the user uuid for the auth request.
func WithUserUUID(userUUID string) AuthOption {
	return func(u *url.Values) {
		u.Set("user_uuid", userUUID)
	}
}
