package groq

// Helper is an interface for error helpers
type Helper interface {
	error
	Advice() string
}
