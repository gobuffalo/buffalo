package binding

import (
	"net/http"
	"time"

	"github.com/gobuffalo/buffalo/binding/decoders"
	"github.com/gobuffalo/nulls"
	"github.com/monoculum/formam"
)

// HTMLContentTypeBinder is in charge of binding HTML request types.
type HTMLContentTypeBinder struct {
	decoder *formam.Decoder
}

// NewHTMLContentTypeBinder returns an instance of HTMLContentTypeBinder with
// custom type decoders registered for Time and nulls.Time
func NewHTMLContentTypeBinder(decoder *formam.Decoder) HTMLContentTypeBinder {
	decoder.RegisterCustomType(decoders.TimeDecoderFn(), []interface{}{time.Time{}}, nil)
	decoder.RegisterCustomType(decoders.NullTimeDecoderFn(), []interface{}{nulls.Time{}}, nil)

	return HTMLContentTypeBinder{
		decoder: decoder,
	}
}

// ContentTypes that will be used to identify HTML requests
func (ht HTMLContentTypeBinder) ContentTypes() []string {
	return []string{
		"application/html",
		"text/html",
		"application/x-www-form-urlencoded",
		"html",
	}
}

// BinderFunc that will take care of the HTML binding
func (ht HTMLContentTypeBinder) BinderFunc() Binder {
	return func(req *http.Request, i interface{}) error {
		err := req.ParseForm()
		if err != nil {
			return err
		}

		if err := ht.decoder.Decode(req.Form, i); err != nil {
			return err
		}

		return nil
	}
}
