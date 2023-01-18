package network

import (
	"context"
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

func (c *Client) Exchange(ctx context.Context, item *state.CollectionItem) (*http.Response, error) {
	uri, err := url.Parse(item.URL)
	if err != nil {
		return nil, err
	}

	headers := item.Headers

	var data io.ReadCloser
	if item.RequestBody != nil {
		data = io.NopCloser(strings.NewReader(item.RequestBody.Payload))
		headers["Content-Type"] = []string{item.RequestBody.ContentType}
	}

	req := &http.Request{
		Method: item.Method,
		URL:    uri,
		Header: item.Headers,
		Body:   data,
	}

	req = req.WithContext(ctx)
	return c.client.Do(req)
}
