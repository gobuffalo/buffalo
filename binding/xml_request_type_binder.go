package binding

import (
	"encoding/xml"
	"net/http"

	"github.com/monoculum/formam"
)

// XMLRequestTypeBinder is in charge of binding XML request types.
type XMLRequestTypeBinder struct{}

func (xm XMLRequestTypeBinder) binder(decoder *formam.Decoder) Binder {
	return func(req *http.Request, value interface{}) error {
		return xml.NewDecoder(req.Body).Decode(value)
	}
}

// RegisterTo register this RequestTypeBinder to the passed request binder
// on the XML content types.
func (xm XMLRequestTypeBinder) RegisterTo(binder *RequestBinder) {
	binder.Register("application/xml", xm.binder(binder.formDecoder))
	binder.Register("text/xml", xm.binder(binder.formDecoder))
	binder.Register("xml", xm.binder(binder.formDecoder))
}
