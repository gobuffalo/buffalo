package binding

import (
	"net/http"

	"github.com/monoculum/formam"
)

type HTMLDecoder struct{}

func (ht HTMLDecoder) Binder(decoder *formam.Decoder) Binder {
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
