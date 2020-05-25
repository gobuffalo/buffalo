package binding

import (
	"encoding/xml"
	"net/http"
)

// XMLRequestTypeBinder is in charge of binding XML request types.
type XMLRequestTypeBinder struct{}

// BinderFunc returns the Binder for this RequestTypeBinder
func (xm XMLRequestTypeBinder) BinderFunc() Binder {
	return func(req *http.Request, value interface{}) error {
		return xml.NewDecoder(req.Body).Decode(value)
	}
}

// ContentTypes that will be wired to this the XML Binder
func (xm XMLRequestTypeBinder) ContentTypes() []string {
	return []string{
		"application/xml",
		"text/xml",
		"xml",
	}
}
