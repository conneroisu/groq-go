# Toolhouse

This is an extension for groq-go that allows you to use the Toolhouse API to give tools to your ai models.

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# toolhouse

```go
import "github.com/conneroisu/groq-go/extensions/toolhouse"
```

Package toolhouse provides a Toolhouse extension for groq\-go.

It allows you to use the Toolhouse API to give tools to your ai models.

Package toolhouse provides a Toolhouse extension for groq\-go.

## Index

- [type Extension](<#Extension>)
  - [func NewExtension\(apiKey string, opts ...Options\) \(e \*Extension, err error\)](<#NewExtension>)
  - [func \(e \*Extension\) GetTools\(ctx context.Context\) \(\[\]groq.Tool, error\)](<#Extension.GetTools>)
  - [func \(e \*Extension\) MustGetTools\(ctx context.Context\) \[\]groq.Tool](<#Extension.MustGetTools>)
  - [func \(e \*Extension\) MustRun\(ctx context.Context, response groq.ChatCompletionResponse\) \[\]groq.ChatCompletionMessage](<#Extension.MustRun>)
  - [func \(e \*Extension\) Run\(ctx context.Context, response groq.ChatCompletionResponse\) \(\[\]groq.ChatCompletionMessage, error\)](<#Extension.Run>)
- [type Options](<#Options>)
  - [func WithBaseURL\(baseURL string\) Options](<#WithBaseURL>)
  - [func WithClient\(client \*http.Client\) Options](<#WithClient>)
  - [func WithLogger\(logger \*slog.Logger\) Options](<#WithLogger>)
  - [func WithMetadata\(metadata map\[string\]any\) Options](<#WithMetadata>)


<a name="Extension"></a>
## type [Extension](<https://github.com/conneroisu/groq-go/blob/main/extensions/toolhouse/toolhouse.go#L24-L34>)

Extension is a Toolhouse extension.

```go
type Extension struct {
    // contains filtered or unexported fields
}
```

<a name="NewExtension"></a>
### func [NewExtension](<https://github.com/conneroisu/groq-go/blob/main/extensions/toolhouse/toolhouse.go#L41>)

```go
func NewExtension(apiKey string, opts ...Options) (e *Extension, err error)
```

NewExtension creates a new Toolhouse extension.

<a name="Extension.GetTools"></a>
### func \(\*Extension\) [GetTools](<https://github.com/conneroisu/groq-go/blob/main/extensions/toolhouse/tools.go#L28-L30>)

```go
func (e *Extension) GetTools(ctx context.Context) ([]groq.Tool, error)
```

GetTools returns a list of tools that the extension can use.

<a name="Extension.MustGetTools"></a>
### func \(\*Extension\) [MustGetTools](<https://github.com/conneroisu/groq-go/blob/main/extensions/toolhouse/tools.go#L17-L19>)

```go
func (e *Extension) MustGetTools(ctx context.Context) []groq.Tool
```

MustGetTools returns a list of tools that the extension can use.

It panics if an error occurs.

<a name="Extension.MustRun"></a>
### func \(\*Extension\) [MustRun](<https://github.com/conneroisu/groq-go/blob/main/extensions/toolhouse/run.go#L26-L29>)

```go
func (e *Extension) MustRun(ctx context.Context, response groq.ChatCompletionResponse) []groq.ChatCompletionMessage
```

MustRun runs the extension on the given history.

It panics if an error occurs.

<a name="Extension.Run"></a>
### func \(\*Extension\) [Run](<https://github.com/conneroisu/groq-go/blob/main/extensions/toolhouse/run.go#L38-L41>)

```go
func (e *Extension) Run(ctx context.Context, response groq.ChatCompletionResponse) ([]groq.ChatCompletionMessage, error)
```

Run runs the extension on the given history.

<a name="Options"></a>
## type [Options](<https://github.com/conneroisu/groq-go/blob/main/extensions/toolhouse/toolhouse.go#L37>)

Options is a function that sets options for a Toolhouse extension.

```go
type Options func(*Extension)
```

<a name="WithBaseURL"></a>
### func [WithBaseURL](<https://github.com/conneroisu/groq-go/blob/main/extensions/toolhouse/options.go#L9>)

```go
func WithBaseURL(baseURL string) Options
```

WithBaseURL sets the base URL for the Toolhouse extension.

<a name="WithClient"></a>
### func [WithClient](<https://github.com/conneroisu/groq-go/blob/main/extensions/toolhouse/options.go#L16>)

```go
func WithClient(client *http.Client) Options
```

WithClient sets the client for the Toolhouse extension.

<a name="WithLogger"></a>
### func [WithLogger](<https://github.com/conneroisu/groq-go/blob/main/extensions/toolhouse/options.go#L30>)

```go
func WithLogger(logger *slog.Logger) Options
```

WithLogger sets the logger for the Toolhouse extension.

<a name="WithMetadata"></a>
### func [WithMetadata](<https://github.com/conneroisu/groq-go/blob/main/extensions/toolhouse/options.go#L23>)

```go
func WithMetadata(metadata map[string]any) Options
```

WithMetadata sets the metadata for the get tools request.

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->
