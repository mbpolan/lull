package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type oauth2Response struct {
	AccessToken string `json:"access_token"`
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

func (a *OAuth2RequestAuthentication) Type() string {
	return "oauth2"
}

func (a *OAuth2RequestAuthentication) Prepare() (*http.Request, error) {
	form := url.Values{}
	form.Set("client_id", a.ClientID)
	form.Set("client_secret", a.ClientSecret)
	form.Set("grant_type", a.GrantType)
	form.Set("scope", a.Scope)

	req, err := http.NewRequest("POST", a.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (a *OAuth2RequestAuthentication) Apply(req *http.Request, res *http.Response) error {
	// validate the status code before attempting to parse the body
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to acquire token: HTTP %d", res.StatusCode)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var oauth2 oauth2Response
	if err := json.Unmarshal(data, &oauth2); err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oauth2.AccessToken))
	return nil
}
