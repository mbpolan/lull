package state

import "net/http"

type AppState struct {
	Method      string
	URL         string
	RequestBody string
	Response    *http.Response
	LastError   error
}
