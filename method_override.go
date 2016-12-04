package buffalo

import (
	"net/http"

	"github.com/markbates/going/defaults"
)

// MethodOverride can be be overridden to a user specified
// function that can be used to change the HTTP Request Method.
// var MethodOverride = MethodOverrideFunc
var MethodOverride http.HandlerFunc

// MethodOverrideFunc is the default implementation for the
// MethodOverride. By default it will look for a form value
// name `_method` and change the request method if that is
// present and the original request is of type "POST". This is
// added automatically when using `Automatic` Buffalo, unless
// an alternative is defined in the Options.
func MethodOverrideFunc(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.Method = defaults.String(req.FormValue("_method"), "POST")
		req.Form.Del("_method")
		req.PostForm.Del("_method")
	}
}
