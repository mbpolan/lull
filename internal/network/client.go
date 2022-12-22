package network

import (
	"net/http"
	"net/url"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	c := new(Client)
	c.client = http.DefaultClient
	return c
}

func (c *Client) Exchange(method string, url *url.URL) (*http.Response, error) {
	return c.client.Do(&http.Request{
		Method: method,
		URL:    url,
	})
}
