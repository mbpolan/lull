package parsers

import (
	"encoding/json"
	"io"
	"net/http"
)

// JSONBodyParser is a parser for "application/json" content types.
type JSONBodyParser struct {
}

// NewJSONBodyParser returns an instance of JSONBodyParser.
func NewJSONBodyParser() *JSONBodyParser {
	return new(JSONBodyParser)
}

// Parse returns a formatted and prettified JSON response body.
func (j *JSONBodyParser) Parse(res *http.Response) (string, error) {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return j.ParseBytes(data)
}

// ParseBytes returns a formatted and prettified JSON response body for the given raw bytes.
func (j *JSONBodyParser) ParseBytes(body []byte) (string, error) {
	var jsonData any
	err := json.Unmarshal(body, &jsonData)
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
