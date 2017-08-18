// +build !appengine

package buffalo

import (
	"context"
	"net/http"
)

func contextFromRequest(req *http.Request) context.Context {
	return req.Context()
}
