package binding

import (
	"encoding/xml"
	"net/http"

	"github.com/monoculum/formam"
)

type XMLDecoder struct{}

func (xm XMLDecoder) Binder(decoder *formam.Decoder) Binder {
	return func(req *http.Request, value interface{}) error {
		return xml.NewDecoder(req.Body).Decode(value)
	}
}
