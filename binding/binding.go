package binding

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/markbates/pop/nulls"
	"github.com/monoculum/formam"
	"github.com/pkg/errors"
)

// Binder takes a request and binds it to an interface.
// If there is a problem it should return an error.
type Binder func(*http.Request, interface{}) error

// CustomTypeDecoder converts a custom type from the request insto its exact type.
type CustomTypeDecoder func([]string) (interface{}, error)

// binders is a map of the defined content-type related binders.
var binders = map[string]Binder{}

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

// RegisterCustomDecorder allows to define custom type decoders.
func RegisterCustomDecorder(fn CustomTypeDecoder, types []interface{}, fields []interface{}) {
	rawFunc := (func([]string) (interface{}, error))(fn)
	decoder.RegisterCustomType(rawFunc, types, fields)
}

// RegisterBinder maps a request Content-Type (application/json)
// to a Binder.
func Register(contentType string, fn Binder) {
	lock.Lock()
	defer lock.Unlock()

	binders[strings.ToLower(contentType)] = fn
}

// Bind the interface to the request.Body. The type of binding
// is dependent on the "Content-Type" for the request. If the type
// is "application/json" it will use "json.NewDecoder". If the type
// is "application/xml" it will use "xml.NewDecoder". The default
// binder is "https://github.com/monoculum/formam".
func Exec(req *http.Request, value interface{}) error {
	ct := strings.ToLower(req.Header.Get("Content-Type"))
	if ct != "" {
		cts := strings.Split(ct, ";")
		c := cts[0]
		if b, ok := binders[strings.TrimSpace(c)]; ok {
			return b(req, value)
		}
		return errors.Errorf("could not find a binder for %s", c)
	}
	return errors.New("blank content type")
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

	decoder.RegisterCustomType(func(vals []string) (interface{}, error) {
		var ti nulls.Time
		var err error

		for _, layout := range timeFormats {
			t, er := time.Parse(layout, vals[0])
			if er == nil {
				ti.Time = t
				return ti, er
			}

			err = er
		}

		return ti, err
	}, []interface{}{nulls.Time{}}, nil)

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

	binders["application/html"] = sb
	binders["text/html"] = sb
	binders["application/x-www-form-urlencoded"] = sb
	binders["multipart/form-data"] = sb
	binders["html"] = sb
}

func init() {
	jb := func(req *http.Request, value interface{}) error {
		return json.NewDecoder(req.Body).Decode(value)
	}

	binders["application/json"] = jb
	binders["text/json"] = jb
	binders["json"] = jb
}

func init() {
	xb := func(req *http.Request, value interface{}) error {
		return xml.NewDecoder(req.Body).Decode(value)
	}

	binders["application/xml"] = xb
	binders["text/xml"] = xb
	binders["xml"] = xb
}
