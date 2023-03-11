package auth

import (
	"net/http"
)

// RequestAuthentication is an interface that allows abstracts authentication scheme implementations.
type RequestAuthentication interface {
	Type() string
	Prepare() (*http.Request, error)
	Apply(req *http.Request, res *http.Response) error
}
