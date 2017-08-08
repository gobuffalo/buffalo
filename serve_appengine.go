// +build appengine

package buffalo

import (
	"net/http"
)

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.MethodOverride != nil {
		a.MethodOverride(w, r)
	}
	a.router.ServeHTTP(w, r)
}
