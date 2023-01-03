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

	var jsonData any
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return "", err
	}

	data, err = json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
