package parsers

import (
	"net/http"
	"strings"
)

const (
	bodyParserNoop string = "noop"
	bodyParserJson string = "json"
)

var instance BodyParserProvider

// BodyParserProvider is a factor that provides instances of various BodyParser implementations.
type BodyParserProvider struct {
	cache map[string]BodyParser
}

// Setup initializes the body parser factory. This should be called once on program start up.
func Setup() {
	instance = BodyParserProvider{}

	noop := NewNoopBodyParser()

	instance.cache = map[string]BodyParser{
		bodyParserNoop: noop,
		bodyParserJson: NewJSONBodyParser(),
	}
}

// GetBodyParser returns an instance of BodyParser most suitable for the HTTP response.
func GetBodyParser(res *http.Response) BodyParser {
	ct := res.Header.Get("Content-Type")
	if ct == "" {
		return instance.cache[bodyParserNoop]
	}

	return GetBodyParserForContentType(ct)
}

// GetBodyParserForContentType returns an instance of BodyParser most suitable for the content type.
func GetBodyParserForContentType(contentType string) BodyParser {
	if strings.Contains(contentType, "application/json") {
		return instance.cache[bodyParserJson]
	}

	return instance.cache[bodyParserNoop]
}
