package binding

import (
	"encoding/json"
	"net/http"
)

// JSONContentTypeBinder is in charge of binding JSON request types.
type JSONContentTypeBinder struct{}

// BinderFunc returns the Binder for this JSONRequestTypeBinder
func (js JSONContentTypeBinder) BinderFunc() Binder {
	return func(req *http.Request, value interface{}) error {
		return json.NewDecoder(req.Body).Decode(value)
	}
}

// ContentTypes that will be wired to this the JSON Binder
func (js JSONContentTypeBinder) ContentTypes() []string {
	return []string{
		"application/json",
		"text/json",
		"json",
	}
}
