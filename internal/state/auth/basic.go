package auth

import "net/http"

// BasicAuthentication contains username and password parameters for HTTP basic auth.
type BasicAuthentication struct {
	Username string
	Password string
}

// NewBasicAuthentication returns a new instance of BasicAuthentication.
func NewBasicAuthentication(username, password string) *BasicAuthentication {
	return &BasicAuthentication{
		Username: username,
		Password: password,
	}
}

func (a *BasicAuthentication) Type() string {
	return "basic"
}

func (a *BasicAuthentication) Prepare() (*http.Request, error) {
	// no additional network request needed
	return nil, nil
}

func (a *BasicAuthentication) Apply(req *http.Request, res *http.Response) error {
	req.SetBasicAuth(a.Username, a.Password)
	return nil
}
