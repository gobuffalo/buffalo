package binding

import (
	"encoding/xml"
	"net/http"

	"github.com/monoculum/formam"
)

type XMLRequestTypeBinder struct{}

func (xm XMLRequestTypeBinder) binder(decoder *formam.Decoder) Binder {
	return func(req *http.Request, value interface{}) error {
		return xml.NewDecoder(req.Body).Decode(value)
	}
}

func (xm XMLRequestTypeBinder) RegisterTo(binder *RequestBinder) {
	binder.Register("application/xml", xm.binder(binder.formDecoder))
	binder.Register("text/xml", xm.binder(binder.formDecoder))
	binder.Register("xml", xm.binder(binder.formDecoder))
}
