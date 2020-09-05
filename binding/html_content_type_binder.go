package binding

import (
	"net/http"

	"github.com/monoculum/formam"
)

// HTMLContentTypeBinder is in charge of binding HTML request types.
type HTMLContentTypeBinder struct {
	decoder *formam.Decoder
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
