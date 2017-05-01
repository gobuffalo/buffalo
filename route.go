package buffalo

import (
	"fmt"
	"strings"

	"reflect"

	"github.com/gorilla/mux"
	"github.com/markbates/inflect"
)

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

	PathName string `json:"pathName"`
	App      *App   `json:"-"`
}

// Name allows users to set custom names for the routes.
func (ri *RouteInfo) Name(name string) *RouteInfo {
	routeIndex := -1
	for index, route := range ri.App.Routes() {
		if route.Path == ri.Path && route.Method == ri.Method {
			routeIndex = index
			break
		}
	}

	name = inflect.CamelizeDownFirst(name)

	if strings.HasSuffix(name, "Path") == false {
		name = name + "Path"
	}

	ri.PathName = name
	if routeIndex != -1 {
		ri.App.Routes()[routeIndex] = reflect.Indirect(reflect.ValueOf(ri)).Interface().(RouteInfo)
	}

	return ri
}

//BuildPathHelper Builds a routeHelperfunc for a particular RouteInfo
func (ri *RouteInfo) BuildPathHelper() RouteHelperFunc {
	cRoute := ri
	return func(opts map[string]interface{}) string {
		pairs := []string{}
		for k, v := range opts {
			pairs = append(pairs, k)
			pairs = append(pairs, fmt.Sprintf("%v", v))
		}

		url, err := cRoute.MuxRoute.URL(pairs...)
		if err != nil {
			return cRoute.Path
		}

		return url.Path
	}
}

//RouteHelperFunc represents the function that takes the route and the opts and build the path
type RouteHelperFunc func(opts map[string]interface{}) string

// RouteList contains a mapping of the routes defined
// in the application. This listing contains, Method, Path,
// and the name of the Handler defined to process that route.
type RouteList []RouteInfo

func (a RouteList) Len() int      { return len(a) }
func (a RouteList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a RouteList) Less(i, j int) bool {
	x := a[i].Path // + a[i].Method
	y := a[j].Path // + a[j].Method
	return x < y
}
