package buffalo

import (
	"net/http"

	"github.com/markbates/going/defaults"
)

var MethodOverride = func(req *http.Request) {
	if req.Method == "POST" {
		req.Method = defaults.String(req.FormValue("_method"), "POST")
	}
}
