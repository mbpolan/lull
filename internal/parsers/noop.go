package parsers

import (
	"io"
	"net/http"
)

// NoopBodyParser is a no-op parser that does not do anything with a response body.
type NoopBodyParser struct {
}

// NewNoopBodyParser returns an instance of NoopBodyParser.
func NewNoopBodyParser() *NoopBodyParser {
	return new(NoopBodyParser)
}

// Parse returns the original HTTP response body as-is.
func (n *NoopBodyParser) Parse(res *http.Response) (string, error) {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ParseBytes returns the original body bytes as-is.
func (n *NoopBodyParser) ParseBytes(body []byte) (string, error) {
	return string(body), nil
}
