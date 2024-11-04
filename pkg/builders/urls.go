package builders

type (
	// URLComputer computes URLs for a given client.
	URLComputer interface {
		ComputeURLs()
	}
)
