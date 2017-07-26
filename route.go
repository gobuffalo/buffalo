package buffalo

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"sort"
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
	PathName    string     `json:"pathName"`
	Aliases     []string   `json:"aliases"`
	MuxRoute    *mux.Route `json:"-"`
	Handler     Handler    `json:"-"`
	App         *App       `json:"-"`
}

// String returns a JSON representation of the RouteInfo
func (ri RouteInfo) String() string {
	b, _ := json.MarshalIndent(ri, "", "  ")
	return string(b)
}

// Alias path patterns to the this route. This is not the
// same as a redirect.
func (ri *RouteInfo) Alias(aliases ...string) *RouteInfo {
	ri.Aliases = append(ri.Aliases, aliases...)
	for _, a := range aliases {
		ri.App.router.Handle(a, ri).Methods(ri.Method)
	}
	return ri
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

	if !strings.HasSuffix(name, "Path") {
		name = name + "Path"
	}

	ri.PathName = name
	if routeIndex != -1 {
		ri.App.Routes()[routeIndex] = reflect.ValueOf(ri).Interface().(*RouteInfo)
	}

	return ri
}

//BuildPathHelper Builds a routeHelperfunc for a particular RouteInfo
func (ri *RouteInfo) BuildPathHelper() RouteHelperFunc {
	cRoute := ri
	return func(opts map[string]interface{}) template.HTML {
		pairs := []string{}
		for k, v := range opts {
			pairs = append(pairs, k)
			pairs = append(pairs, fmt.Sprintf("%v", v))
		}

		url, err := cRoute.MuxRoute.URL(pairs...)
		if err != nil {
			return template.HTML(cRoute.Path)
		}

		result := url.Path
		result = addExtraParamsTo(result, opts)

		return template.HTML(result)
	}
}

func addExtraParamsTo(path string, opts map[string]interface{}) string {
	pendingParams := map[string]string{}
	keys := []string{}
	for k, v := range opts {
		if strings.Contains(path, fmt.Sprintf("%v", v)) {
			continue
		}

		keys = append(keys, k)
		pendingParams[k] = fmt.Sprintf("%v", v)
	}

	if len(keys) == 0 {
		return path
	}

	if strings.Contains(path, "?") == false {
		path = path + "?"
	} else {
		if strings.HasSuffix(path, "?") == false {
			path = path + "&"
		}
	}

	sort.Strings(keys)

	for index, k := range keys {
		format := "%v=%v"

		if index > 0 {
			format = "&%v=%v"
		}

		path = path + fmt.Sprintf(format, url.QueryEscape(k), url.QueryEscape(pendingParams[k]))
		index = index + 1
	}

	return path
}

//RouteHelperFunc represents the function that takes the route and the opts and build the path
type RouteHelperFunc func(opts map[string]interface{}) template.HTML

// RouteList contains a mapping of the routes defined
// in the application. This listing contains, Method, Path,
// and the name of the Handler defined to process that route.
type RouteList []*RouteInfo

func (a RouteList) Len() int      { return len(a) }
func (a RouteList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a RouteList) Less(i, j int) bool {
	x := a[i].Path // + a[i].Method
	y := a[j].Path // + a[j].Method
	return x < y
}
