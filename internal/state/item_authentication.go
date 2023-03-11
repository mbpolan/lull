package state

import (
	"encoding/json"
	"github.com/mbpolan/lull/internal/state/auth"
)

// ItemAuthentication is a container for authentication parameters for a collection item.
type ItemAuthentication struct {
	Data auth.RequestAuthentication
}

func (i *ItemAuthentication) None() bool {
	return i.Data == nil
}

func (i *ItemAuthentication) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)

	if i.Data == nil {
		m["Type"] = "none"
	} else {
		m["Type"] = i.Data.Type()
		m["Data"] = i.Data
	}

	return json.Marshal(m)
}

func (i *ItemAuthentication) UnmarshalJSON(b []byte) error {
	var raw map[string]*json.RawMessage
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	var authType string
	err = json.Unmarshal(*raw["Type"], &authType)
	if err != nil {
		return err
	}

	if authType == (&auth.OAuth2RequestAuthentication{}).Type() {
		var oauth2 auth.OAuth2RequestAuthentication
		err = json.Unmarshal(*raw["Data"], &oauth2)
		if err != nil {
			return err
		}

		i.Data = &oauth2
	} else if authType == (&auth.BasicAuthentication{}).Type() {
		var basic auth.BasicAuthentication
		err = json.Unmarshal(*raw["Data"], &basic)
		if err != nil {
			return err
		}

		i.Data = &basic
	}

	return nil
}
