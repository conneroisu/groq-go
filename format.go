package groq

// Format is the format of a response.
// string
type Format string

const (
	// FormatText is the text format. It is the default format of a
	// response.
	FormatText Format = "text"
	// FormatJSON is the JSON format. There is no support for streaming with
	// JSON format selected.
	FormatJSON Format = "json"
)
