package binding

import (
	"encoding/json"
	"net/http"

	"github.com/monoculum/formam"
)

type JSONDecoder struct{}

func (js JSONDecoder) Binder(decoder *formam.Decoder) Binder {
	return func(req *http.Request, value interface{}) error {
		return json.NewDecoder(req.Body).Decode(value)
	}
}
