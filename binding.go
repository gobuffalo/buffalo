package buffalo

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"sync"
	"time"

	"github.com/monoculum/formam"
	"github.com/pkg/errors"
)

var binderLock = &sync.Mutex{}
var binders = map[string]BinderFunc{}
var decoder *formam.Decoder

var timeLayouts = []string{
	"2006-01-02T15:04:05Z07:00",
}

// BinderFunc takes a request and binds it to an interface.
// If there is a problem it should return an error.
type BinderFunc func(*http.Request, interface{}) error

// RegisterBinder maps a request Content-Type (application/json)
// to a BinderFunc.
func RegisterBinder(contentType string, fn BinderFunc) {
	binderLock.Lock()
	defer binderLock.Unlock()
	binders[contentType] = fn
}

// RegisterTimeLayout allows to add custom time layouts that
// the binder will be able to use for decoding.
func RegisterTimeLayout(layout string) {
	timeLayouts = append(timeLayouts, layout)
}

// RegisterCustomTypeDecoder
func RegisterBinderTypeDecoder(fn formam.DecodeCustomTypeFunc, types []interface{}, fields []interface{}) {
	decoder.RegisterCustomType(fn, types, fields)
}

func init() {
	decoder = formam.NewDecoder(&formam.DecoderOptions{
		TagName:           "schema",
		IgnoreUnknownKeys: true,
	})

	decoder.RegisterCustomType(func(vals []string) (interface{}, error) {
		var t time.Time
		var err error

		for _, layout := range timeLayouts {
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

	binders["application/html"] = sb
	binders["text/html"] = sb
	binders["application/x-www-form-urlencoded"] = sb
	binders["multipart/form-data"] = sb
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
