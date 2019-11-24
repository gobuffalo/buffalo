package buffalo

import (
	"net/http"
)

// MethodOverride is the default implementation for the
// Options#MethodOverride. By default it will look for a form value
// name `_method` and change the request method if that is
// present and the original request is of type "POST". This is
// added automatically when using `New` Buffalo, unless
// an alternative is defined in the Options.
func MethodOverride(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		md := req.FormValue("_method")
		if len(md) != 0 {
			req.Method = md
		}
		req.Form.Del("_method")
		req.PostForm.Del("_method")
	}
}
