package binding

import (
	"net/http"
	"time"

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

	// timeFormats are the base time formats supported by the time.Time and
	// nulls.Time Decoders you can prepend custom formats to this list
	// by using RegisterTimeFormats.
	timeFormats = []string{
		time.RFC3339,
		"01/02/2006",
		"2006-01-02",
		"2006-01-02T15:04",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	// BaseRequestBinder is an instance of the requeBinder, it comes with preconfigured
	// content type binders for HTML, JSON, XML and Files, as well as custom types decoders
	// for time.Time and nulls.Time
	BaseRequestBinder = NewRequestBinder(
		HTMLContentTypeBinder{
			decoder: formDecoder,
		},
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
	timeFormats = append(layouts, timeFormats...)
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
