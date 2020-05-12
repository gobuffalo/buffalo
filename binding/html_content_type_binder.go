package binding

import (
	"net/http"

	"github.com/monoculum/formam"
)

// HTMLContentTypeBinder is in charge of binding HTML request types.
type HTMLContentTypeBinder struct {
	decoder *formam.Decoder
}

func NewHTMLContentTypeBinder(decoder *formam.Decoder) HTMLContentTypeBinder {
	htmlBinder := HTMLContentTypeBinder{
		decoder: decoder,
	}

	// decoder.RegisterCustomType(, types []interface{}, fields []interface{})

	return htmlBinder
}

func (ht HTMLContentTypeBinder) ContentTypes() []string {
	return []string{
		"application/html",
		"text/html",
		"application/x-www-form-urlencoded",
		"html",
	}
}

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
