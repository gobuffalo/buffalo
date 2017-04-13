package buffalo

import (
	"fmt"
	"net/http"
	"path"
	"reflect"
	"sort"
	"strings"

	"github.com/markbates/inflect"
)

// GET maps an HTTP "GET" request to the path and the specified handler.
func (a *App) GET(p string, h Handler) *RouteInfo {
	return a.addRoute("GET", p, h)
}

// POST maps an HTTP "POST" request to the path and the specified handler.
func (a *App) POST(p string, h Handler) *RouteInfo {
	return a.addRoute("POST", p, h)
}

// PUT maps an HTTP "PUT" request to the path and the specified handler.
func (a *App) PUT(p string, h Handler) *RouteInfo {
	return a.addRoute("PUT", p, h)
}

// DELETE maps an HTTP "DELETE" request to the path and the specified handler.
func (a *App) DELETE(p string, h Handler) *RouteInfo {
	return a.addRoute("DELETE", p, h)
}

// HEAD maps an HTTP "HEAD" request to the path and the specified handler.
func (a *App) HEAD(p string, h Handler) *RouteInfo {
	return a.addRoute("HEAD", p, h)
}

// OPTIONS maps an HTTP "OPTIONS" request to the path and the specified handler.
func (a *App) OPTIONS(p string, h Handler) *RouteInfo {
	return a.addRoute("OPTIONS", p, h)
}

// PATCH maps an HTTP "PATCH" request to the path and the specified handler.
func (a *App) PATCH(p string, h Handler) *RouteInfo {
	return a.addRoute("PATCH", p, h)
}

// Redirect from one URL to another URL. Only works for "GET" requests.
func (a *App) Redirect(status int, from, to string) *RouteInfo {
	return a.GET(from, func(c Context) error {
		return c.Redirect(status, to)
	})
}

// ServeFiles maps an path to a directory on disk to serve static files.
// Useful for JavaScript, images, CSS, etc...
/*
	a.ServeFiles("/assets", http.Dir("path/to/assets"))
*/
func (a *App) ServeFiles(p string, root http.FileSystem) {
	a.router.PathPrefix(p).Handler(http.StripPrefix(p, http.FileServer(root)))
}

// Resource maps an implementation of the Resource interface
// to the appropriate RESTful mappings. Resource returns the *App
// associated with this group of mappings so you can set middleware, etc...
// on that group, just as if you had used the a.Group functionality.
/*
	a.Resource("/users", &UsersResource{})

	// Is equal to this:

	ur := &UsersResource{}
	g := a.Group("/users")
	g.GET("/", ur.List) // GET /users => ur.List
	g.GET("/new", ur.New) // GET /users/new => ur.New
	g.GET("/{user_id}", ur.Show) // GET /users/{user_id} => ur.Show
	g.GET("/{user_id}/edit", ur.Edit) // GET /users/{user_id}/edit => ur.Edit
	g.POST("/", ur.Create) // POST /users => ur.Create
	g.PUT("/{user_id}", ur.Update) PUT /users/{user_id} => ur.Update
	g.DELETE("/{user_id}", ur.Destroy) DELETE /users/{user_id} => ur.Destroy
*/
func (a *App) Resource(p string, r Resource) *App {
	base := path.Base(p)
	single := inflect.Singularize(base)
	g := a.Group(p)
	p = "/"

	rv := reflect.ValueOf(r)
	rt := rv.Type()
	rname := fmt.Sprintf("%s.%s", rt.PkgPath(), rt.Name()) + ".%s"

	spath := path.Join(p, fmt.Sprintf("{%s_id}", single))
	setFuncKey(r.List, fmt.Sprintf(rname, "List"))
	g.GET(p, r.List)
	setFuncKey(r.New, fmt.Sprintf(rname, "New"))
	g.GET(path.Join(p, "new"), r.New).Name(inflect.Camelize(fmt.Sprintf("new_" + single)))
	setFuncKey(r.Show, fmt.Sprintf(rname, "Show"))
	g.GET(path.Join(spath), r.Show)
	setFuncKey(r.Edit, fmt.Sprintf(rname, "Edit"))
	g.GET(path.Join(spath, "edit"), r.Edit).Name(inflect.Camelize(fmt.Sprintf("edit_" + single)))
	setFuncKey(r.Create, fmt.Sprintf(rname, "Create"))
	g.POST(p, r.Create)
	setFuncKey(r.Update, fmt.Sprintf(rname, "Update"))
	g.PUT(path.Join(spath), r.Update)
	setFuncKey(r.Destroy, fmt.Sprintf(rname, "Destroy"))
	g.DELETE(path.Join(spath), r.Destroy)
	return g
}

// ANY accepts a request across any HTTP method for the specified path
// and routes it to the specified Handler.
func (a *App) ANY(p string, h Handler) {
	a.GET(p, h)
	a.POST(p, h)
	a.PUT(p, h)
	a.PATCH(p, h)
	a.HEAD(p, h)
	a.OPTIONS(p, h)
	a.DELETE(p, h)
}

// Group creates a new `*App` that inherits from it's parent `*App`.
// This is useful for creating groups of end-points that need to share
// common functionality, like middleware.
/*
	g := a.Group("/api/v1")
	g.Use(AuthorizeAPIMiddleware)
	g.GET("/users, APIUsersHandler)
	g.GET("/users/:user_id, APIUserShowHandler)
*/
func (a *App) Group(groupPath string) *App {
	g := New(a.Options)

	g.prefix = path.Join(a.prefix, groupPath)

	g.router = a.router
	g.Middleware = a.Middleware.clone()
	g.ErrorHandlers = a.ErrorHandlers
	g.root = a
	if a.root != nil {
		g.root = a.root
	}
	return g
}

func (a *App) addRoute(method string, url string, h Handler) *RouteInfo {
	a.moot.Lock()
	defer a.moot.Unlock()

	url = path.Join(a.prefix, url)

	hs := funcKey(h)
	r := RouteInfo{
		Method:      method,
		Path:        url,
		HandlerName: hs,
		Handler:     h,
		App:         a,
	}

	r.MuxRoute = a.router.Handle(url, a.handlerToHandler(r, h)).Methods(method)
	r.Name(buildRouteName(url))

	routes := a.Routes()
	routes = append(routes, r)
	sort.Sort(routes)

	if a.root != nil {
		a.root.routes = routes
	} else {
		a.routes = routes
	}

	return &r
}

//buildRouteName builds a route based on the path passed.
func buildRouteName(path string) string {

	if path == "/" {
		return "root"
	}

	resultPars := []string{}
	parts := strings.Split(path, "/")

	for index, part := range parts {

		if strings.Contains(part, "{") || part == "" {
			continue
		}
		shouldSingularize := (len(parts) > index+1) && strings.Contains(parts[index+1], "{")
		if shouldSingularize {
			part = inflect.Singularize(part)
		}

		if index > 0 && strings.Contains(parts[index-1], "}") {
			resultPars = append(resultPars, part)
			continue
		}

		resultPars = append([]string{part}, resultPars...)
	}

	underscore := strings.Join(resultPars, "_")
	return inflect.CamelizeDownFirst(underscore)
}
