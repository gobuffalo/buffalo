// +build appengine

package buffalo

import (
	"context"
	"net/http"

	"google.golang.org/appengine"
)

func contextFromRequest(req *http.Request) context.Context {
	return appengine.NewContext(req)
}
