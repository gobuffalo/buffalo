package buffalo

import (
	"net/http"

	"github.com/gobuffalo/buffalo/internal/defaults"
)

// MethodOverride is the default implementation for the
// Options#MethodOverride. By default it will look for a form value
// name `_method` and change the request method if that is
// present and the original request is of type "POST". This is
// added automatically when using `New` Buffalo, unless
// an alternative is defined in the Options.
func MethodOverride(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.Method = defaults.String(req.FormValue("_method"), "POST")
		req.Form.Del("_method")
		req.PostForm.Del("_method")
	}
}
