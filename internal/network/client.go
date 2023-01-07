package network

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	c := new(Client)
	c.client = http.DefaultClient
	return c
}

func (c *Client) Exchange(method string, url *url.URL, body string, headers map[string][]string) (*http.Response, error) {
	var data io.ReadCloser
	if body != "" {
		data = io.NopCloser(strings.NewReader(body))
	}

	return c.client.Do(&http.Request{
		Method: method,
		URL:    url,
		Header: headers,
		Body:   data,
	})
}
