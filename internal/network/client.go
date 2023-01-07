package network

import (
	"github.com/mbpolan/lull/internal/state"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Client struct {
	client *http.Client
	mutex  sync.Mutex
}

func NewClient() *Client {
	c := new(Client)
	c.client = http.DefaultClient

	return c
}

func (c *Client) Exchange(item *state.CollectionItem) (*http.Response, error) {
	uri, err := url.Parse(item.URL)
	if err != nil {
		return nil, err
	}

	var data io.ReadCloser
	if item.RequestBody != "" {
		data = io.NopCloser(strings.NewReader(item.RequestBody))
	}

	return c.client.Do(&http.Request{
		Method: item.Method,
		URL:    uri,
		Header: item.Headers,
		Body:   data,
	})
}
