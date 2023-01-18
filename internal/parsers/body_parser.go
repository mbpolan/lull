package parsers

import "net/http"

// BodyParser is the top-level interface for implementations that parse HTTP response bodies.
type BodyParser interface {
	// Parse reads and formats the HTTP response body returning a formatted string.
	Parse(res *http.Response) (string, error)

	// ParseBytes formats a raw byte representation of body content, returning a formatted string.
	ParseBytes(body []byte) (string, error)
}
