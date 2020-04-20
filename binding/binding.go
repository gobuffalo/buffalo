package binding

import (
	"net/http"
)

// MaxFileMemory can be used to set the maximum size, in bytes, for files to be
// stored in memory during uploaded for multipart requests.
// See https://golang.org/pkg/net/http/#Request.ParseMultipartForm for more
// information on how this impacts file uploads.
var MaxFileMemory int64 = 5 * 1024 * 1024

// RequestBinder is an instance of the default request binder, it comes with preconfigured
// content type binders for HTML, JSON, XML and Files, as well as custom types decoders
// for time.Time and nulls.Time
var defaultRequestBinder = NewDefaultRequestBinder(
	HTMLRequestTypeBinder{},
	JSONRequestTypeBinder{},
	XMLRequestTypeBinder{},
	FileRequestTypeBinder{},
)

// Binder takes a request and binds it to an interface.
// If there is a problem it should return an error.
type Binder func(*http.Request, interface{}) error

// CustomTypeDecoder converts a custom type from the request insto its exact type.
type CustomTypeDecoder func([]string) (interface{}, error)

// RequestTypeBinder are those capable of handling a request type like JSON or XML
type RequestTypeBinder interface {
	RegisterTo(*RequestBinder)
}

// RegisterTimeFormats allows to add custom time layouts that
// the binder will be able to use for decoding.
func RegisterTimeFormats(layouts ...string) {
	defaultRequestBinder.timeFormats = append(layouts, defaultRequestBinder.timeFormats...)
}

// RegisterCustomDecoder allows to define custom type decoders.
func RegisterCustomDecoder(fn CustomTypeDecoder, types []interface{}, fields []interface{}) {
	defaultRequestBinder.RegisterCustomDecoder(fn, types, fields)
}

// Register maps a request Content-Type (application/json)
// to a Binder.
func Register(contentType string, fn Binder) {
	defaultRequestBinder.Register(contentType, fn)
}

// Exec will bind the interface to the request.Body. The type of binding
// is dependent on the "Content-Type" for the request. If the type
// is "application/json" it will use "json.NewDecoder". If the type
// is "application/xml" it will use "xml.NewDecoder". The default
// binder is "https://github.com/monoculum/formam".
func Exec(req *http.Request, value interface{}) error {
	return defaultRequestBinder.Exec(req, value)
}
