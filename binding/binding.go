package binding

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"sync"
	"time"

	"github.com/monoculum/formam"
	"github.com/pkg/errors"
)

// BinderFunc takes a request and binds it to an interface.
// If there is a problem it should return an error.
type BinderFunc func(*http.Request, interface{}) error

// CustomTypeDecoderFunc converts a custom type from the request insto its exact type.
type CustomTypeDecoderFunc func([]string) (interface{}, error)

// Binders is a map of the defined content-type related binders.
var Binders = map[string]BinderFunc{}

var decoder *formam.Decoder
var lock = &sync.Mutex{}
var timeFormats = []string{
	"2006-01-02T15:04:05Z07:00",
}

// RegisterTimeFormats allows to add custom time layouts that
// the binder will be able to use for decoding.
func RegisterTimeFormats(layouts ...string) {
	timeFormats = append(timeFormats, layouts...)
}

// RegisterBinderTypeDecoder allows to define custom type decoders.
func RegisterBinderTypeDecoder(fn CustomTypeDecoderFunc, types []interface{}, fields []interface{}) {
	rawFunc := (func([]string) (interface{}, error))(fn)
	decoder.RegisterCustomType(rawFunc, types, fields)
}

// RegisterBinder maps a request Content-Type (application/json)
// to a BinderFunc.
func RegisterBinder(contentType string, fn BinderFunc) {
	lock.Lock()
	defer lock.Unlock()

	Binders[contentType] = fn
}

func init() {
	decoder = formam.NewDecoder(&formam.DecoderOptions{
		TagName:           "form",
		IgnoreUnknownKeys: true,
	})

	decoder.RegisterCustomType(func(vals []string) (interface{}, error) {
		var t time.Time
		var err error

		for _, layout := range timeFormats {
			t, er := time.Parse(layout, vals[0])
			if er == nil {
				return t, er
			}

			err = er
		}

		return t, err
	}, []interface{}{time.Time{}}, nil)

	sb := func(req *http.Request, i interface{}) error {
		err := req.ParseForm()
		if err != nil {
			return errors.WithStack(err)
		}

		if err := decoder.Decode(req.Form, i); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}

	Binders["application/html"] = sb
	Binders["text/html"] = sb
	Binders["application/x-www-form-urlencoded"] = sb
	Binders["multipart/form-data"] = sb
}

func init() {
	jb := func(req *http.Request, value interface{}) error {
		return json.NewDecoder(req.Body).Decode(value)
	}

	Binders["application/json"] = jb
	Binders["text/json"] = jb
	Binders["json"] = jb
}

func init() {
	xb := func(req *http.Request, value interface{}) error {
		return xml.NewDecoder(req.Body).Decode(value)
	}

	Binders["application/xml"] = xb
	Binders["text/xml"] = xb
	Binders["xml"] = xb
}
