package binding

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"errors"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/x/httpx"
	"github.com/monoculum/formam"
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

// RegisterTimeFormats allows to add custom time layouts that
// the binder will be able to use for decoding.
func RegisterTimeFormats(layouts ...string) {
	timeFormats = append(layouts, timeFormats...)
}

// RegisterCustomDecoder allows to define custom type decoders.
func RegisterCustomDecoder(fn CustomTypeDecoder, types []interface{}, fields []interface{}) {
	rawFunc := (func([]string) (interface{}, error))(fn)
	decoder.RegisterCustomType(rawFunc, types, fields)
}

// Register maps a request Content-Type (application/json)
// to a Binder.
func Register(contentType string, fn Binder) {
	lock.Lock()
	defer lock.Unlock()

	binders[strings.ToLower(contentType)] = fn
}

// Exec will bind the interface to the request.Body. The type of binding
// is dependent on the "Content-Type" for the request. If the type
// is "application/json" it will use "json.NewDecoder". If the type
// is "application/xml" it will use "xml.NewDecoder". The default
// binder is "https://github.com/monoculum/formam".
func Exec(req *http.Request, value interface{}) error {
	if ba, ok := value.(Bindable); ok {
		return ba.Bind(req)
	}

	ct := httpx.ContentType(req)
	if ct == "" {
		return errors.New("blank content type")
	}
	if b, ok := binders[ct]; ok {
		return b(req, value)
	}
	return fmt.Errorf("could not find a binder for %s", ct)
}

func init() {
	decoder = formam.NewDecoder(&formam.DecoderOptions{
		TagName:           "form",
		IgnoreUnknownKeys: true,
	})

	decoder.RegisterCustomType(func(vals []string) (interface{}, error) {
		return parseTime(vals)
	}, []interface{}{time.Time{}}, nil)

	decoder.RegisterCustomType(func(vals []string) (interface{}, error) {
		var ti nulls.Time

		t, err := parseTime(vals)
		if err != nil {
			return ti, err
		}
		ti.Time = t
		ti.Valid = true

		return ti, nil
	}, []interface{}{nulls.Time{}}, nil)

	sb := func(req *http.Request, i interface{}) error {
		err := req.ParseForm()
		if err != nil {
			return err
		}

		if err := decoder.Decode(req.Form, i); err != nil {
			return err
		}
		return nil
	}

	binders["application/html"] = sb
	binders["text/html"] = sb
	binders["application/x-www-form-urlencoded"] = sb
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

func parseTime(vals []string) (time.Time, error) {
	var t time.Time
	var err error

	// don't try to parse empty time values, it will raise an error
	if len(vals) == 0 || vals[0] == "" {
		return t, nil
	}

	for _, layout := range timeFormats {
		t, err = time.Parse(layout, vals[0])
		if err == nil {
			return t, nil
		}
	}

	if err != nil {
		return t, err
	}

	return t, nil
}
