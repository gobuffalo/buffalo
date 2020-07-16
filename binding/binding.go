package binding

import (
	"net/http"

	"github.com/gobuffalo/buffalo/binding/decoders"
	"github.com/monoculum/formam"
)

var (
	// MaxFileMemory can be used to set the maximum size, in bytes, for files to be
	// stored in memory during uploaded for multipart requests.
	// See https://golang.org/pkg/net/http/#Request.ParseMultipartForm for more
	// information on how this impacts file uploads.
	MaxFileMemory int64 = 5 * 1024 * 1024

	formDecoder = formam.NewDecoder(&formam.DecoderOptions{
		TagName:           "form",
		IgnoreUnknownKeys: true,
	})

	// BaseRequestBinder is an instance of the requestBinder, it comes with preconfigured
	// content type binders for HTML, JSON, XML and Files, as well as custom types decoders
	// for time.Time and nulls.Time
	BaseRequestBinder = NewRequestBinder(
		NewHTMLContentTypeBinder(formDecoder),
		JSONContentTypeBinder{},
		XMLRequestTypeBinder{},
		FileRequestTypeBinder{
			decoder: formDecoder,
		},
	)
)

// RegisterTimeFormats allows to add custom time layouts that
// the binder will be able to use for decoding.
func RegisterTimeFormats(layouts ...string) {
	decoders.RegisterTimeFormats(layouts...)
}

// RegisterCustomDecoder allows to define custom decoders for certain types
// In the request.
func RegisterCustomDecoder(fn CustomTypeDecoder, types []interface{}, fields []interface{}) {
	rawFunc := (func([]string) (interface{}, error))(fn)
	formDecoder.RegisterCustomType(rawFunc, types, fields)
}

// Register maps a request Content-Type (application/json)
// to a Binder.
func Register(contentType string, fn Binder) {
	BaseRequestBinder.Register(contentType, fn)
}

// Exec will bind the interface to the request.Body. The type of binding
// is dependent on the "Content-Type" for the request. If the type
// is "application/json" it will use "json.NewDecoder". If the type
// is "application/xml" it will use "xml.NewDecoder". The default
// binder is "https://github.com/monoculum/formam".
func Exec(req *http.Request, value interface{}) error {
	return BaseRequestBinder.Exec(req, value)
}
