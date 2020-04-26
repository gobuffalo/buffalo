package binding

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gobuffalo/buffalo/internal/httpx"
	"github.com/gobuffalo/nulls"
	"github.com/monoculum/formam"
)

var (
	errBlankContentType = errors.New("blank content type")
)

// RequestBinder is in charge of binding multiple requests types to
// structs.
type RequestBinder struct {
	lock    *sync.Mutex
	binders map[string]Binder

	formDecoder *formam.Decoder
	timeFormats []string
}

// Register maps a request Content-Type (application/json)
// to a Binder.
func (rb RequestBinder) Register(contentType string, fn Binder) {
	rb.lock.Lock()
	defer rb.lock.Unlock()

	rb.binders[strings.ToLower(contentType)] = fn
}

// RegisterCustomDecoder allows to define custom decoders for certain types
// In the request.
func (rb RequestBinder) RegisterCustomDecoder(fn CustomTypeDecoder, types []interface{}, fields []interface{}) {
	rawFunc := (func([]string) (interface{}, error))(fn)
	rb.formDecoder.RegisterCustomType(rawFunc, types, fields)
}

// Exec binds a request with a passed value, depending on the content type
// It will look for the correct RequestTypeBinder and use it.
func (rb RequestBinder) Exec(req *http.Request, value interface{}) error {
	if ba, ok := value.(Bindable); ok {
		return ba.Bind(req)
	}

	ct := httpx.ContentType(req)
	if ct == "" {
		return errBlankContentType
	}

	if b, ok := rb.binders[ct]; ok {
		return b(req, value)
	}

	return fmt.Errorf("could not find a binder for %s", ct)
}

// NewRequestBinder creates our request binder with support for
// XML, JSON, HTTP and File request types. It also adds decoders
// for Time and nulls.Time.
func NewRequestBinder(requestBinders ...RequestTypeBinder) *RequestBinder {
	result := &RequestBinder{
		lock:    &sync.Mutex{},
		binders: map[string]Binder{},

		formDecoder: formam.NewDecoder(&formam.DecoderOptions{
			TagName:           "form",
			IgnoreUnknownKeys: true,
		}),

		timeFormats: []string{
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
		},
	}

	timeCustom := TimeCustomTypeDecoder{&result.timeFormats}
	result.formDecoder.RegisterCustomType(timeCustom.Decode, []interface{}{time.Time{}}, nil)

	nullTimeCustom := NullTimeCustomTypeDecoder{&timeCustom}
	result.formDecoder.RegisterCustomType(nullTimeCustom.Decode, []interface{}{nulls.Time{}}, nil)

	for _, requestBinder := range requestBinders {
		requestBinder.RegisterTo(result)
	}

	return result
}
