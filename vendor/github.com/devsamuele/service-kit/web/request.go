package web

import (
	"encoding/json"
	"net/http"

	"github.com/dimfeld/httptreemux"
)

// URIParams ...
func URIParams(r *http.Request) map[string]string {
	return httptreemux.ContextParams(r.Context())
}

// QueryParams ...
func QueryParams(r *http.Request) map[string]string {
	query := r.URL.Query()
	result := make(map[string]string)
	for k, v := range query {
		result[k] = v[0]
	}

	return result
}

// Decode ...
func Decode(r *http.Request, data interface{}) error {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}
