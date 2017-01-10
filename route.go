package buffalo

import "github.com/gorilla/mux"

// Routes returns a list of all of the routes defined
// in this application.
func (a *App) Routes() RouteList {
	if a.root != nil {
		return a.root.routes
	}
	return a.routes
}

// RouteInfo provides information about the underlying route that
// was built.
type RouteInfo struct {
	Method      string     `json:"method"`
	Path        string     `json:"path"`
	HandlerName string     `json:"handler"`
	MuxRoute    *mux.Route `json:"-"`
	Handler     Handler    `json:"-"`
}

// RouteList contains a mapping of the routes defined
// in the application. This listing contains, Method, Path,
// and the name of the Handler defined to process that route.
type RouteList []RouteInfo

func (a RouteList) Len() int      { return len(a) }
func (a RouteList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a RouteList) Less(i, j int) bool {
	x := a[i].Method + a[i].Path
	y := a[j].Method + a[j].Path
	return x < y
}
