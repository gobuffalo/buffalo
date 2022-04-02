package buffalo

import (
	"fmt"
	"html/template"
	"net/url"
	"sort"
	"strings"
)

// Routes returns a list of all of the routes defined
// in this application.
func (a *App) Routes() RouteList {
	// CHKME: why this function is exported? can we deprecate it?
	if a.root != nil {
		return a.root.routes
	}
	return a.routes
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

	if !strings.Contains(path, "?") {
		path = path + "?"
	} else {
		if !strings.HasSuffix(path, "?") {
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
	}

	return path
}

//RouteHelperFunc represents the function that takes the route and the opts and build the path
type RouteHelperFunc func(opts map[string]interface{}) (template.HTML, error)

// RouteList contains a mapping of the routes defined
// in the application. This listing contains, Method, Path,
// and the name of the Handler defined to process that route.
type RouteList []*RouteInfo

var methodOrder = map[string]string{
	"GET":    "1",
	"POST":   "2",
	"PUT":    "3",
	"DELETE": "4",
}

func (a RouteList) Len() int      { return len(a) }
func (a RouteList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a RouteList) Less(i, j int) bool {
	// NOTE: it was used for sorting of app.routes but we don't sort the routes anymore.
	// keep it for compatibility but could be deprecated.
	x := a[i].App.host + a[i].Path + methodOrder[a[i].Method]
	y := a[j].App.host + a[j].Path + methodOrder[a[j].Method]
	return x < y
}

// Lookup search a specific PathName in the RouteList and return the *RouteInfo
func (a RouteList) Lookup(name string) (*RouteInfo, error) {
	for _, ri := range a {
		if ri.PathName == name {
			return ri, nil
		}
	}
	return nil, fmt.Errorf("path name not found")
}
