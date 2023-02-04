package state

// RequestAuthentication is an interface that allows abstracts authentication scheme implementations.
type RequestAuthentication interface {
}

// OAuth2RequestAuthentication contains authentication parameters for an OAuth2 token.
type OAuth2RequestAuthentication struct {
	TokenURL     string
	ClientID     string
	ClientSecret string
	GrantType    string
	Scope        string
}

// NewOAuth2RequestAuthentication returns a new instance of OAuth2RequestAuthentication.
func NewOAuth2RequestAuthentication(tokenURL, clientID, clientSecret, grantType, scope string) *OAuth2RequestAuthentication {
	return &OAuth2RequestAuthentication{
		TokenURL:     tokenURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    grantType,
		Scope:        scope,
	}
}
