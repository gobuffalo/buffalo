package binding

import (
	"net/http"

	"github.com/monoculum/formam"
)

// HTMLRequestTypeBinder is in charge of binding HTML request types.
type HTMLRequestTypeBinder struct{}

// RegisterTo register this RequestTypeBinder to the passed request binder
// on the HTML content types.
func (ht HTMLRequestTypeBinder) RegisterTo(binder *RequestBinder) {
	binder.Register("application/html", ht.binder(binder.formDecoder))
	binder.Register("text/html", ht.binder(binder.formDecoder))
	binder.Register("application/x-www-form-urlencoded", ht.binder(binder.formDecoder))
	binder.Register("html", ht.binder(binder.formDecoder))
}

func (ht HTMLRequestTypeBinder) binder(decoder *formam.Decoder) Binder {
	return func(req *http.Request, i interface{}) error {
		err := req.ParseForm()
		if err != nil {
			return err
		}

		if err := decoder.Decode(req.Form, i); err != nil {
			return err
		}
		return nil
	}
}
