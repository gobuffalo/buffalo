package binding

import (
	"encoding/json"
	"net/http"

	"github.com/monoculum/formam"
)

type JSONRequestTypeBinder struct{}

func (js JSONRequestTypeBinder) binder(decoder *formam.Decoder) Binder {
	return func(req *http.Request, value interface{}) error {
		return json.NewDecoder(req.Body).Decode(value)
	}
}

func (js JSONRequestTypeBinder) RegisterTo(binder *RequestBinder) {
	binder.Register("application/json", js.binder(binder.formDecoder))
	binder.Register("text/json", js.binder(binder.formDecoder))
	binder.Register("json", js.binder(binder.formDecoder))
}
