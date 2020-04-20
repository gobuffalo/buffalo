package binding

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gobuffalo/buffalo/internal/httpx"
	"github.com/gobuffalo/nulls"
	"github.com/monoculum/formam"
)

type DefaultRequestBinder struct {
	lock        *sync.Mutex
	binders     map[string]Binder
	formDecoder *formam.Decoder
	timeFormats []string
}

// Register maps a request Content-Type (application/json)
// to a Binder.
func (rb DefaultRequestBinder) Register(contentType string, fn Binder) {
	rb.lock.Lock()
	defer rb.lock.Unlock()

	rb.binders[strings.ToLower(contentType)] = fn
}

func (rb DefaultRequestBinder) RegisterCustomDecoder(fn CustomTypeDecoder, types []interface{}, fields []interface{}) {
	rawFunc := (func([]string) (interface{}, error))(fn)
	rb.formDecoder.RegisterCustomType(rawFunc, types, fields)
}

func (rb DefaultRequestBinder) Exec(req *http.Request, value interface{}) error {
	if ba, ok := value.(Bindable); ok {
		return ba.Bind(req)
	}

	ct := httpx.ContentType(req)
	if ct == "" {
		return fmt.Errorf("blank content type")
	}

	if b, ok := rb.binders[ct]; ok {
		return b(req, value)
	}

	return fmt.Errorf("could not find a binder for %s", ct)
}

func NewDefaultRequestBinder() *DefaultRequestBinder {
	result := &DefaultRequestBinder{
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

	htmlDecoder := HTMLDecoder{}
	result.Register("application/html", htmlDecoder.Binder(result.formDecoder))
	result.Register("text/html", htmlDecoder.Binder(result.formDecoder))
	result.Register("application/x-www-form-urlencoded", htmlDecoder.Binder(result.formDecoder))
	result.Register("html", htmlDecoder.Binder(result.formDecoder))

	xmlDecoder := XMLDecoder{}
	result.Register("application/xml", xmlDecoder.Binder(result.formDecoder))
	result.Register("text/xml", xmlDecoder.Binder(result.formDecoder))
	result.Register("xml", xmlDecoder.Binder(result.formDecoder))

	jsonDecoder := JSONDecoder{}
	result.Register("application/json", jsonDecoder.Binder(result.formDecoder))
	result.Register("text/json", jsonDecoder.Binder(result.formDecoder))
	result.Register("json", jsonDecoder.Binder(result.formDecoder))

	fileDecoder := FileDecoder{}
	result.Register("multipart/form-data", fileDecoder.Binder(result.formDecoder))

	return result
}
